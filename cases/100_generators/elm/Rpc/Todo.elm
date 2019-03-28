module Rpc.Todo exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E
import RpcUtil exposing (Config, CallResult, unwrapHttpResult, decodeCallResult)



type alias User =
    { ctime : ()
    , username : String
    }

encodeUser : User -> E.Value
encodeUser obj =
    E.object
        [ ( "ctime", never obj.ctime )
        , ( "username", E.string obj.username )
        ]

decodeUser : D.Decoder User
decodeUser =
    D.map2 User
            (D.field "ctime" (never))
            (D.field "username" (D.string))
    




