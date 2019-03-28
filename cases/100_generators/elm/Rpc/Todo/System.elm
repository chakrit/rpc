module Rpc.Todo.System exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E
import RpcUtil exposing (Config, CallResult, unwrapHttpResult, decodeCallResult)
import Rpc





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
