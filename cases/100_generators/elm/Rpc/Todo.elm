module Rpc.Todo exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E
import Time
import RpcUtil exposing (Config, CallResult, unwrapHttpResult, decodeCallResult)



type alias User =
    { ctime : Time.Posix
    , username : String
    }

encodeUser : User -> E.Value
encodeUser obj =
    E.object
        [ ( "ctime", (Time.posixToMillis >> toFloat >> (\f -> f/1000.0) >> E.float) obj.ctime )
        , ( "username", E.string obj.username )
        ]

decodeUser : D.Decoder User
decodeUser =
    D.map2 User
            (D.field "ctime" ((D.map ((\f -> f * 1000.0) >> round >> Time.millisToPosix) D.float)))
            (D.field "username" (D.string))
    




