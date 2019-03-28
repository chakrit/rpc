module Rpc exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E
import RpcUtil exposing (Config, CallResult, unwrapHttpResult, decodeCallResult)



type alias Failure =
    { code : String
    , description : String
    }

encodeFailure : Failure -> E.Value
encodeFailure obj =
    E.object
        [ ( "code", E.string obj.code )
        , ( "description", E.string obj.description )
        ]

decodeFailure : D.Decoder Failure
decodeFailure =
    D.map2 Failure
            (D.field "code" (D.string))
            (D.field "description" (D.string))
    




