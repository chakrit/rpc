module RpcUtil exposing (..)

import Array exposing (Array)
import Http exposing (Error(..))
import Json.Decode as D


type alias Config =
    { baseUrl : String
    , headers : List Http.Header
    }


decodeConfig : D.Decoder Config
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
    D.map2 Config
        (D.field "baseUrl" <| D.string)
        (D.field "headers" <| D.map mapHeaderArray (D.array (D.array D.string)))


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


unwrapHttpResult : (Result Http.Error (CallResult a)) -> CallResult a
unwrapHttpResult result =
    case result of
        Ok callResult ->
            callResult

        Err httpErr ->
            HttpError httpErr


decodeCallResult : D.Decoder a -> D.Decoder (CallResult a)
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
    D.map2 (mapResultObj)
        (D.field "error" (D.maybe D.string))
        (D.field "returns" decodeReturns)
