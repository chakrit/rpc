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
        [ ( "ctime", (Time.posixToMillis >> E.int) obj.ctime )
        , ( "username", E.string obj.username )
        ]

decodeUser : D.Decoder User
decodeUser =
    D.map2 User
            (D.field "ctime" ((D.map Time.millisToPosix D.int)))
            (D.field "username" (D.string))
    




