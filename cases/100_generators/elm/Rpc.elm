module Rpc exposing (..)

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
    




