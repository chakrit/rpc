module RpcUtil exposing
    ( CallResult(..)
    , Config
    , base64FromBytes
    , base64ToBytes
    , decodeCallResult
    , decodeField
    , decodeConfig
    , fromBase64String
    , mapResult
    , toBase64String
    , translateHttpError
    , unwrapHttpResult
    )

import Array exposing (Array)
import Bitwise as Bits
import Bytes exposing (Bytes)
import Bytes.Decode as BytesDec
import Bytes.Encode as BytesEnc
import Dict exposing (Dict)
import Http exposing (Error(..))
import Json.Decode as JsonDec


type alias Config =
    { baseUrl : String
    , headers : List Http.Header
    }


decodeConfig : JsonDec.Decoder Config
decodeConfig =
    let
        mapHeader : Array String -> Http.Header
        mapHeader arr =
            let
                header =
                    Array.get 0 arr |> Maybe.withDefault ""

                content =
                    Array.get 1 arr |> Maybe.withDefault ""
            in
            Http.header header content

        mapHeaderArray : Array (Array String) -> List Http.Header
        mapHeaderArray arr =
            arr |> Array.map mapHeader |> Array.toList
    in
    JsonDec.map2 Config
        (JsonDec.field "baseUrl" <| JsonDec.string)
        (JsonDec.field "headers" <| JsonDec.map mapHeaderArray (JsonDec.array (JsonDec.array JsonDec.string)))

decodeField : String -> JsonDec.Decoder a -> JsonDec.Decoder (a -> b) -> JsonDec.Decoder b
decodeField key fieldDec partial =
    JsonDec.andThen (\p -> JsonDec.map p (JsonDec.field key fieldDec)) partial

type CallResult a
    = HttpError Http.Error
    | ApiError String
    | Success a


translateHttpError : Http.Error -> String
translateHttpError httpErr =
    case httpErr of
        BadUrl str ->
            "Bad URL: " ++ str

        Timeout ->
            "Network Timeout"

        NetworkError ->
            "Network Error"

        BadStatus code ->
            "Bad Status Code: " ++ String.fromInt code

        BadBody str ->
            "Malformed Response: " ++ str


mapResult : { onHttpErr : Http.Error -> msg, onApiErr : String -> msg, onSuccess : a -> msg } -> CallResult a -> msg
mapResult mapper result =
    case result of
        HttpError httpErr ->
            mapper.onHttpErr httpErr

        ApiError err ->
            mapper.onApiErr err

        Success obj ->
            mapper.onSuccess obj


unwrapHttpResult : Result Http.Error (CallResult a) -> CallResult a
unwrapHttpResult result =
    case result of
        Ok callResult ->
            callResult

        Err httpErr ->
            HttpError httpErr


decodeCallResult : JsonDec.Decoder a -> JsonDec.Decoder (CallResult a)
decodeCallResult decodeReturns =
    let
        mapResultObj : Maybe String -> a -> CallResult a
        mapResultObj err ret =
            case err of
                Just str ->
                    ApiError str

                Nothing ->
                    Success ret
    in
    JsonDec.map2 mapResultObj
        (JsonDec.field "error" (JsonDec.maybe JsonDec.string))
        (JsonDec.field "returns" decodeReturns)


toBase64String : String -> Maybe String
toBase64String =
    stringToBytes >> base64FromBytes


fromBase64String : String -> Maybe String
fromBase64String =
    base64ToBytes >> Maybe.map stringFromBytes


stringToBytes : String -> Bytes
stringToBytes str =
    BytesEnc.encode (BytesEnc.string str)


stringFromBytes : Bytes -> String
stringFromBytes bytes =
    BytesDec.decode (BytesDec.string (Bytes.width bytes)) bytes
        |> Maybe.withDefault "(string decode error)"


base64ToBytes : String -> Maybe Bytes
base64ToBytes str =
    decode (String.toList str)
        |> Maybe.map (BytesEnc.sequence >> BytesEnc.encode)


decode : List Char -> Maybe (List BytesEnc.Encoder)
decode l =
    case l of
        [] ->
            Just []

        a :: b :: '=' :: '=' :: [] ->
            dec2 a b |> Maybe.map List.singleton

        a :: b :: [] ->
            dec2 a b |> Maybe.map List.singleton

        a :: b :: c :: '=' :: [] ->
            dec3 a b c |> Maybe.map List.singleton

        a :: b :: c :: [] ->
            dec3 a b c |> Maybe.map List.singleton

        a :: b :: c :: d :: tail ->
            case ( dec4 a b c d, decode tail ) of
                ( Just e, Just list ) ->
                    Just <| e :: list

                _ ->
                    Nothing

        _ ->
            Nothing


{-| base64 1111 1122 2222 3333 3344 4444
|||| bytes 1111 1111 2222 2222 3333 3333
-}
dec4 : Char -> Char -> Char -> Char -> Maybe BytesEnc.Encoder
dec4 a b c d =
    case [ b64byte a, b64byte b, b64byte c, b64byte d ] of
        [ Just aa, Just bb, Just cc, Just dd ] ->
            Just <|
                BytesEnc.sequence
                    [ BytesEnc.unsignedInt8 ((aa |> Bits.shiftLeftBy 2) + (bb |> Bits.shiftRightZfBy 4))
                    , BytesEnc.unsignedInt8 ((bb |> Bits.shiftLeftBy 4) + (cc |> Bits.shiftRightZfBy 2))
                    , BytesEnc.unsignedInt8 ((cc |> Bits.shiftLeftBy 6) + dd)
                    ]

        _ ->
            Nothing


{-| base64 1111 1122 2222 3333 33
|||| bytes 1111 1111 2222 2222 ..
-}
dec3 : Char -> Char -> Char -> Maybe BytesEnc.Encoder
dec3 a b c =
    case ( b64byte a, b64byte b, b64byte c ) of
        ( Just aa, Just bb, Just cc ) ->
            Just <|
                BytesEnc.sequence
                    [ BytesEnc.unsignedInt8 ((aa |> Bits.shiftLeftBy 2) + (bb |> Bits.shiftRightZfBy 4))
                    , BytesEnc.unsignedInt8 ((bb |> Bits.shiftLeftBy 4) + (cc |> Bits.shiftRightZfBy 2))
                    ]

        _ ->
            Nothing


{-| base64 1111 1122 2222 ....
|||| bytes 1111 1111 .... ....
-}
dec2 : Char -> Char -> Maybe BytesEnc.Encoder
dec2 a b =
    case ( b64byte a, b64byte b ) of
        ( Just aa, Just bb ) ->
            Just <| BytesEnc.unsignedInt8 ((aa |> Bits.shiftLeftBy 2) + (bb |> Bits.shiftRightZfBy 4))

        _ ->
            Nothing


base64FromBytes : Bytes -> Maybe String
base64FromBytes bytes =
    let
        get3 =
            BytesDec.map3 enc3 BytesDec.unsignedInt8 BytesDec.unsignedInt8 BytesDec.unsignedInt8

        get2 =
            BytesDec.map2 enc2 BytesDec.unsignedInt8 BytesDec.unsignedInt8

        get1 =
            BytesDec.map enc1 BytesDec.unsignedInt8

        process : DecoderState -> BytesDec.Decoder (BytesDec.Step DecoderState (List Char))
        process state =
            if state.processed >= state.total then
                BytesDec.succeed <| BytesDec.Done state.output

            else
                BytesDec.map BytesDec.Loop <|
                    case state.total - state.processed of
                        1 ->
                            get1 |> BytesDec.map (appendState 1 state)

                        2 ->
                            get2 |> BytesDec.map (appendState 2 state)

                        _ ->
                            -- 3 or more
                            get3 |> BytesDec.map (appendState 3 state)

        decoder =
            BytesDec.loop (initialState bytes) process
                |> BytesDec.map String.fromList
    in
    BytesDec.decode decoder bytes


{-| bytes 1111 1111 2222 2222 3333 3333
|| base64 1111 1122 2222 3333 3344 4444
-}
enc3 : Int -> Int -> Int -> List Char
enc3 a b c =
    [ Bits.shiftRightBy 2 a |> Bits.and 0x3F |> b64char
    , (Bits.shiftLeftBy 4 a + Bits.shiftRightBy 4 b) |> Bits.and 0x3F |> b64char
    , (Bits.shiftLeftBy 2 b + Bits.shiftRightBy 6 c) |> Bits.and 0x3F |> b64char
    , c |> Bits.and 0x3F |> b64char
    ]


{-| bytes 1111 1111 2222 2222 .... ....
|| base64 1111 1122 2222 3333 33.. ....
-}
enc2 : Int -> Int -> List Char
enc2 a b =
    [ Bits.shiftRightBy 2 a |> Bits.and 0x3F |> b64char
    , (Bits.shiftLeftBy 4 a + Bits.shiftRightBy 4 b) |> Bits.and 0x3F |> b64char
    , Bits.shiftLeftBy 2 b |> Bits.and 0x3F |> b64char
    , '='
    ]


{-| bytes 1111 1111 .... ....
|| base64 1111 1122 2222 ....
-}
enc1 : Int -> List Char
enc1 a =
    [ Bits.shiftRightBy 2 a |> Bits.and 0x3F |> b64char
    , Bits.shiftLeftBy 4 a |> Bits.and 0x3F |> b64char
    , '='
    , '='
    ]


type alias DecoderState =
    { total : Int
    , processed : Int
    , output : List Char
    }


initialState : Bytes -> DecoderState
initialState bytes =
    { total = Bytes.width bytes
    , processed = 0
    , output = []
    }


appendState : Int -> DecoderState -> List Char -> DecoderState
appendState count state chars =
    { state
        | processed = state.processed + count
        , output = List.append state.output chars
    }


b64byte : Char -> Maybe Int
b64byte c =
    Dict.get c backwardTable


b64char : Int -> Char
b64char n =
    Dict.get n forwardTable
        |> Maybe.withDefault '?'


backwardTable : Dict Char Int
backwardTable =
    Dict.fromList
        [ ( 'A', 0 )
        , ( 'B', 1 )
        , ( 'C', 2 )
        , ( 'D', 3 )
        , ( 'E', 4 )
        , ( 'F', 5 )
        , ( 'G', 6 )
        , ( 'H', 7 )
        , ( 'I', 8 )
        , ( 'J', 9 )
        , ( 'K', 10 )
        , ( 'L', 11 )
        , ( 'M', 12 )
        , ( 'N', 13 )
        , ( 'O', 14 )
        , ( 'P', 15 )
        , ( 'Q', 16 )
        , ( 'R', 17 )
        , ( 'S', 18 )
        , ( 'T', 19 )
        , ( 'U', 20 )
        , ( 'V', 21 )
        , ( 'W', 22 )
        , ( 'X', 23 )
        , ( 'Y', 24 )
        , ( 'Z', 25 )
        , ( 'a', 26 )
        , ( 'b', 27 )
        , ( 'c', 28 )
        , ( 'd', 29 )
        , ( 'e', 30 )
        , ( 'f', 31 )
        , ( 'g', 32 )
        , ( 'h', 33 )
        , ( 'i', 34 )
        , ( 'j', 35 )
        , ( 'k', 36 )
        , ( 'l', 37 )
        , ( 'm', 38 )
        , ( 'n', 39 )
        , ( 'o', 40 )
        , ( 'p', 41 )
        , ( 'q', 42 )
        , ( 'r', 43 )
        , ( 's', 44 )
        , ( 't', 45 )
        , ( 'u', 46 )
        , ( 'v', 47 )
        , ( 'w', 48 )
        , ( 'x', 49 )
        , ( 'y', 50 )
        , ( 'z', 51 )
        , ( '0', 52 )
        , ( '1', 53 )
        , ( '2', 54 )
        , ( '3', 55 )
        , ( '4', 56 )
        , ( '5', 57 )
        , ( '6', 58 )
        , ( '7', 59 )
        , ( '8', 60 )
        , ( '9', 61 )
        , ( '-', 62 )
        , ( '_', 63 )
        ]


forwardTable : Dict Int Char
forwardTable =
    Dict.fromList
        [ ( 0, 'A' )
        , ( 1, 'B' )
        , ( 2, 'C' )
        , ( 3, 'D' )
        , ( 4, 'E' )
        , ( 5, 'F' )
        , ( 6, 'G' )
        , ( 7, 'H' )
        , ( 8, 'I' )
        , ( 9, 'J' )
        , ( 10, 'K' )
        , ( 11, 'L' )
        , ( 12, 'M' )
        , ( 13, 'N' )
        , ( 14, 'O' )
        , ( 15, 'P' )
        , ( 16, 'Q' )
        , ( 17, 'R' )
        , ( 18, 'S' )
        , ( 19, 'T' )
        , ( 20, 'U' )
        , ( 21, 'V' )
        , ( 22, 'W' )
        , ( 23, 'X' )
        , ( 24, 'Y' )
        , ( 25, 'Z' )
        , ( 26, 'a' )
        , ( 27, 'b' )
        , ( 28, 'c' )
        , ( 29, 'd' )
        , ( 30, 'e' )
        , ( 31, 'f' )
        , ( 32, 'g' )
        , ( 33, 'h' )
        , ( 34, 'i' )
        , ( 35, 'j' )
        , ( 36, 'k' )
        , ( 37, 'l' )
        , ( 38, 'm' )
        , ( 39, 'n' )
        , ( 40, 'o' )
        , ( 41, 'p' )
        , ( 42, 'q' )
        , ( 43, 'r' )
        , ( 44, 's' )
        , ( 45, 't' )
        , ( 46, 'u' )
        , ( 47, 'v' )
        , ( 48, 'w' )
        , ( 49, 'x' )
        , ( 50, 'y' )
        , ( 51, 'z' )
        , ( 52, '0' )
        , ( 53, '1' )
        , ( 54, '2' )
        , ( 55, '3' )
        , ( 56, '4' )
        , ( 57, '5' )
        , ( 58, '6' )
        , ( 59, '7' )
        , ( 60, '8' )
        , ( 61, '9' )
        , ( 62, '-' )
        , ( 63, '_' )
        ]

