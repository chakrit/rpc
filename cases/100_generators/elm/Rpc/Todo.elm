module Rpc.Todo exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E

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
    




