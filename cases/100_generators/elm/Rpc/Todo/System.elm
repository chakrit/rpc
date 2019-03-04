module Rpc.Todo.System exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E
import Rpc

type alias Config =
    { baseUrl : String
    , headers : List Http.Header
    }


type CallResult a
    = HttpError Http.Error
    | ApiError String
    | Success a


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





type alias InputForStatus =
    (())

encodeInputForStatus : InputForStatus -> E.Value
encodeInputForStatus
    () =
        E.list (identity)
            [
            ]

decodeInputForStatus : D.Decoder InputForStatus
decodeInputForStatus =
        D.succeed ()

type alias OutputForStatus =
    (Rpc.Failure)

encodeOutputForStatus : OutputForStatus -> E.Value
encodeOutputForStatus
    (arg0) =
        E.list (identity)
            [ Rpc.encodeFailure arg0
            ]

decodeOutputForStatus : D.Decoder OutputForStatus
decodeOutputForStatus =
        D.map (\a -> (a))
            (D.index 0 (Rpc.decodeFailure))



callStatus : Config -> InputForStatus -> (CallResult OutputForStatus -> a) -> Cmd a
callStatus config input mapResult =
    Http.request
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/rpc/todo/system/Status"
        , body = Http.jsonBody (encodeInputForStatus input)
        , expect = Http.expectJson (unwrapHttpResult >> mapResult) (decodeCallResult (decodeOutputForStatus))
        , timeout = Nothing
        , tracker = Nothing
        }
