module Rpc.Todo.Auth exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E
import RpcUtil exposing (Config, CallResult, unwrapHttpResult, decodeCallResult)
import Rpc
import Rpc.Todo



type alias AuthRequest =
    { provider : String
    , secret : String
    , username : String
    }

encodeAuthRequest : AuthRequest -> E.Value
encodeAuthRequest obj =
    E.object
        [ ( "provider", E.string obj.provider )
        , ( "secret", E.string obj.secret )
        , ( "username", E.string obj.username )
        ]

decodeAuthRequest : D.Decoder AuthRequest
decodeAuthRequest =
    D.map3 AuthRequest
            (D.field "provider" (D.string))
            (D.field "secret" (D.string))
            (D.field "username" (D.string))
    

type alias AuthResponse =
    { failure : Rpc.Failure
    , user : Rpc.Todo.User
    }

encodeAuthResponse : AuthResponse -> E.Value
encodeAuthResponse obj =
    E.object
        [ ( "failure", Rpc.encodeFailure obj.failure )
        , ( "user", Rpc.Todo.encodeUser obj.user )
        ]

decodeAuthResponse : D.Decoder AuthResponse
decodeAuthResponse =
    D.map2 AuthResponse
            (D.field "failure" (Rpc.decodeFailure))
            (D.field "user" (Rpc.Todo.decodeUser))
    



type alias InputForAuthenticate =
    (AuthRequest)

encodeInputForAuthenticate : InputForAuthenticate -> E.Value
encodeInputForAuthenticate
    (arg0) =
        E.list (identity)
            [ encodeAuthRequest arg0
            ]

decodeInputForAuthenticate : D.Decoder InputForAuthenticate
decodeInputForAuthenticate =
        D.map (\a -> (a))
            (D.index 0 (decodeAuthRequest))

type alias OutputForAuthenticate =
    (AuthResponse)

encodeOutputForAuthenticate : OutputForAuthenticate -> E.Value
encodeOutputForAuthenticate
    (arg0) =
        E.list (identity)
            [ encodeAuthResponse arg0
            ]

decodeOutputForAuthenticate : D.Decoder OutputForAuthenticate
decodeOutputForAuthenticate =
        D.map (\a -> (a))
            (D.index 0 (decodeAuthResponse))

type alias InputForCurrent =
    (())

encodeInputForCurrent : InputForCurrent -> E.Value
encodeInputForCurrent
    () =
        E.list (identity)
            [
            ]

decodeInputForCurrent : D.Decoder InputForCurrent
decodeInputForCurrent =
        D.succeed ()

type alias OutputForCurrent =
    (Rpc.Todo.User)

encodeOutputForCurrent : OutputForCurrent -> E.Value
encodeOutputForCurrent
    (arg0) =
        E.list (identity)
            [ Rpc.Todo.encodeUser arg0
            ]

decodeOutputForCurrent : D.Decoder OutputForCurrent
decodeOutputForCurrent =
        D.map (\a -> (a))
            (D.index 0 (Rpc.Todo.decodeUser))



callAuthenticate : Config -> InputForAuthenticate -> (CallResult OutputForAuthenticate -> a) -> Cmd a
callAuthenticate config input mapResult =
    Http.request
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/rpc/todo/auth/Authenticate"
        , body = Http.jsonBody (encodeInputForAuthenticate input)
        , expect = Http.expectJson (unwrapHttpResult >> mapResult) (decodeCallResult (decodeOutputForAuthenticate))
        , timeout = Nothing
        , tracker = Nothing
        }

callCurrent : Config -> InputForCurrent -> (CallResult OutputForCurrent -> a) -> Cmd a
callCurrent config input mapResult =
    Http.request
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/rpc/todo/auth/Current"
        , body = Http.jsonBody (encodeInputForCurrent input)
        , expect = Http.expectJson (unwrapHttpResult >> mapResult) (decodeCallResult (decodeOutputForCurrent))
        , timeout = Nothing
        , tracker = Nothing
        }
