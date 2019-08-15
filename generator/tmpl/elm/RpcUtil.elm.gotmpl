module RpcUtil exposing
    ( Config
    , RpcError(..)
    , RpcResult
    , b64StringFromBytes
    , b64StringToBytes
    , configDecoder
    , decodeApply
    , decodeString
    , decodeValue
    , decoder
    , emptyBytes
    , errorToString
    , fromHttpResult
    , map
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


configDecoder : JsonDec.Decoder Config
configDecoder =
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


decodeApply : JsonDec.Decoder a -> JsonDec.Decoder (a -> b) -> JsonDec.Decoder b
decodeApply fieldDec partial =
    JsonDec.andThen (\p -> JsonDec.map p fieldDec) partial


type RpcError
    = HttpError Http.Error
    | JsonError JsonDec.Error
    | ApiError String


type alias RpcResult a =
    Result RpcError a


fromHttpResult : Result Http.Error (RpcResult a) -> RpcResult a
fromHttpResult httpResult =
    case httpResult of
        Err err ->
            Err (HttpError err)

        Ok result ->
            -- unwrap inner result
            result


errorToString : RpcError -> String
errorToString httpErr =
    case httpErr of
        HttpError (BadUrl str) ->
            "Bad URL: " ++ str

        HttpError Timeout ->
            "Network Timeout"

        HttpError NetworkError ->
            "Network Error"

        HttpError (BadStatus code) ->
            "Bad Status Code: " ++ String.fromInt code

        HttpError (BadBody str) ->
            "Malformed Response: " ++ str

        JsonError err ->
            "JSON Error: " ++ JsonDec.errorToString err

        ApiError str ->
            str


map : (RpcError -> msg) -> (a -> msg) -> RpcResult a -> msg
map errMap okMap result =
    case result of
        Err err ->
            errMap err

        Ok obj ->
            okMap obj


decodeString : JsonDec.Decoder (RpcResult a) -> String -> RpcResult a
decodeString decoder_ str =
    case JsonDec.decodeString decoder_ str of
        Ok v ->
            -- unwrap inner Result
            v

        Err err ->
            Err (JsonError err)


decodeValue : JsonDec.Decoder (RpcResult a) -> JsonDec.Value -> RpcResult a
decodeValue decoder_ value =
    case JsonDec.decodeValue decoder_ value of
        Ok v ->
            -- unwrap inner Result
            v

        Err err ->
            Err (JsonError err)


decoder : JsonDec.Decoder a -> JsonDec.Decoder (RpcResult a)
decoder returnDecoder =
    let
        mapToResult : Maybe String -> a -> RpcResult a
        mapToResult err ret =
            case err of
                Just str ->
                    Err (ApiError str)

                Nothing ->
                    Ok ret
    in
    JsonDec.map2 mapToResult
        (JsonDec.field "error" (JsonDec.maybe JsonDec.string))
        (JsonDec.field "returns" returnDecoder)


emptyBytes : Bytes
emptyBytes =
    BytesEnc.encode (BytesEnc.sequence [])


stringToBytes : String -> Bytes
stringToBytes str =
    BytesEnc.encode (BytesEnc.string str)


stringFromBytes : Bytes -> Maybe String
stringFromBytes bytes =
    BytesDec.decode (BytesDec.string (Bytes.width bytes)) bytes


b64StringToBytes : String -> Maybe Bytes
b64StringToBytes str =
    decodeB64Chars (String.toList str)
        |> Maybe.map (BytesEnc.sequence >> BytesEnc.encode)


decodeB64Chars : List Char -> Maybe (List BytesEnc.Encoder)
decodeB64Chars l =
    case l of
        [] ->
            Just []

        a :: b :: '=' :: '=' :: [] ->
            b64Dec2 a b |> Maybe.map List.singleton

        a :: b :: [] ->
            b64Dec2 a b |> Maybe.map List.singleton

        a :: b :: c :: '=' :: [] ->
            b64Dec3 a b c |> Maybe.map List.singleton

        a :: b :: c :: [] ->
            b64Dec3 a b c |> Maybe.map List.singleton

        a :: b :: c :: d :: tail ->
            case ( b64Dec4 a b c d, decodeB64Chars tail ) of
                ( Just e, Just list ) ->
                    Just <| e :: list

                _ ->
                    Nothing

        _ ->
            Nothing


b64StringFromBytes : Bytes -> Maybe String
b64StringFromBytes bytes =
    let
        get3 =
            BytesDec.map3 b64Enc3 BytesDec.unsignedInt8 BytesDec.unsignedInt8 BytesDec.unsignedInt8

        get2 =
            BytesDec.map2 b64Enc2 BytesDec.unsignedInt8 BytesDec.unsignedInt8

        get1 =
            BytesDec.map b64Enc1 BytesDec.unsignedInt8

        process : B64State -> BytesDec.Decoder (BytesDec.Step B64State (List Char))
        process state =
            if state.processed >= state.total then
                BytesDec.succeed <| BytesDec.Done state.output

            else
                BytesDec.map BytesDec.Loop <|
                    case state.total - state.processed of
                        1 ->
                            get1 |> BytesDec.map (appendB64State 1 state)

                        2 ->
                            get2 |> BytesDec.map (appendB64State 2 state)

                        _ ->
                            -- 3 or more
                            get3 |> BytesDec.map (appendB64State 3 state)

        decoder_ =
            BytesDec.loop (initialB64State bytes) process
                |> BytesDec.map String.fromList
    in
    BytesDec.decode decoder_ bytes


{-| base64 1111 1122 2222 3333 3344 4444
|||| bytes 1111 1111 2222 2222 3333 3333
-}
b64Dec4 : Char -> Char -> Char -> Char -> Maybe BytesEnc.Encoder
b64Dec4 a b c d =
    case [ toB64byte a, toB64byte b, toB64byte c, toB64byte d ] of
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
b64Dec3 : Char -> Char -> Char -> Maybe BytesEnc.Encoder
b64Dec3 a b c =
    case ( toB64byte a, toB64byte b, toB64byte c ) of
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
b64Dec2 : Char -> Char -> Maybe BytesEnc.Encoder
b64Dec2 a b =
    case ( toB64byte a, toB64byte b ) of
        ( Just aa, Just bb ) ->
            Just <| BytesEnc.unsignedInt8 ((aa |> Bits.shiftLeftBy 2) + (bb |> Bits.shiftRightZfBy 4))

        _ ->
            Nothing


{-| bytes 1111 1111 2222 2222 3333 3333
|| base64 1111 1122 2222 3333 3344 4444
-}
b64Enc3 : Int -> Int -> Int -> List Char
b64Enc3 a b c =
    [ Bits.shiftRightBy 2 a |> Bits.and 0x3F |> toB64char
    , (Bits.shiftLeftBy 4 a + Bits.shiftRightBy 4 b) |> Bits.and 0x3F |> toB64char
    , (Bits.shiftLeftBy 2 b + Bits.shiftRightBy 6 c) |> Bits.and 0x3F |> toB64char
    , c |> Bits.and 0x3F |> toB64char
    ]


{-| bytes 1111 1111 2222 2222 .... ....
|| base64 1111 1122 2222 3333 33.. ....
-}
b64Enc2 : Int -> Int -> List Char
b64Enc2 a b =
    [ Bits.shiftRightBy 2 a |> Bits.and 0x3F |> toB64char
    , (Bits.shiftLeftBy 4 a + Bits.shiftRightBy 4 b) |> Bits.and 0x3F |> toB64char
    , Bits.shiftLeftBy 2 b |> Bits.and 0x3F |> toB64char
    , '='
    ]


{-| bytes 1111 1111 .... ....
|| base64 1111 1122 2222 ....
-}
b64Enc1 : Int -> List Char
b64Enc1 a =
    [ Bits.shiftRightBy 2 a |> Bits.and 0x3F |> toB64char
    , Bits.shiftLeftBy 4 a |> Bits.and 0x3F |> toB64char
    , '='
    , '='
    ]


type alias B64State =
    { total : Int
    , processed : Int
    , output : List Char
    }


initialB64State : Bytes -> B64State
initialB64State bytes =
    { total = Bytes.width bytes
    , processed = 0
    , output = []
    }


appendB64State : Int -> B64State -> List Char -> B64State
appendB64State count state chars =
    { state
        | processed = state.processed + count
        , output = List.append state.output chars
    }


toB64byte : Char -> Maybe Int
toB64byte c =
    Dict.get c b64BackwardTable


toB64char : Int -> Char
toB64char n =
    Dict.get n b64ForwardTable
        |> Maybe.withDefault '?'


b64BackwardTable : Dict Char Int
b64BackwardTable =
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


b64ForwardTable : Dict Int Char
b64ForwardTable =
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
