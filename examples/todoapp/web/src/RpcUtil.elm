module RpcUtil exposing
    ( Config
    , RpcError(..)
    , RpcResult
    , configDecoder
    , decodeApply
    , decodeString
    , decodeValue
    , decoder
    , errorToString
    , fromHttpResult
    , map
    , resolver
    )

import Array exposing (Array)
import Bitwise as Bits
import Bytes exposing (Bytes)
import Bytes.Decode as BytesDec
import Bytes.Encode as BytesEnc
import Dict exposing (Dict)
import Http exposing (Error(..), Resolver, Response(..))
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


resolver : JsonDec.Decoder a -> Resolver RpcError a
resolver decoder_ =
    Http.stringResolver
        (\resp ->
            case resp of
                BadUrl_ s ->
                    Err (HttpError (BadUrl s))

                Timeout_ ->
                    Err (HttpError Timeout)

                NetworkError_ ->
                    Err (HttpError NetworkError)

                BadStatus_ metadata _ ->
                    Err (HttpError (BadStatus metadata.statusCode))

                GoodStatus_ metadata str ->
                    decodeString (decoder decoder_) str
        )


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

