module Api exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E
import RpcUtil exposing (Config, CallResult, unwrapHttpResult, decodeCallResult)



type alias TodoItem =
    { description : String
    , done : Bool
    , id : String
    }

encodeTodoItem : TodoItem -> E.Value
encodeTodoItem obj =
    E.object
        [ ( "description", E.string obj.description )
        , ( "done", E.bool obj.done )
        , ( "id", E.string obj.id )
        ]

decodeTodoItem : D.Decoder TodoItem
decodeTodoItem =
    D.map3 TodoItem
            (D.field "description" (D.string))
            (D.field "done" (D.bool))
            (D.field "id" (D.string))
    



type alias InputForDestroy =
    (String)

encodeInputForDestroy : InputForDestroy -> E.Value
encodeInputForDestroy
    (arg0) =
        E.list (identity)
            [ E.string arg0
            ]

decodeInputForDestroy : D.Decoder InputForDestroy
decodeInputForDestroy =
        D.map (\a -> (a))
            (D.index 0 (D.string))

type alias OutputForDestroy =
    (TodoItem)

encodeOutputForDestroy : OutputForDestroy -> E.Value
encodeOutputForDestroy
    (arg0) =
        E.list (identity)
            [ encodeTodoItem arg0
            ]

decodeOutputForDestroy : D.Decoder OutputForDestroy
decodeOutputForDestroy =
        D.map (\a -> (a))
            (D.index 0 (decodeTodoItem))

type alias InputForList =
    (())

encodeInputForList : InputForList -> E.Value
encodeInputForList
    () =
        E.list (identity)
            [
            ]

decodeInputForList : D.Decoder InputForList
decodeInputForList =
        D.succeed ()

type alias OutputForList =
    (List (TodoItem))

encodeOutputForList : OutputForList -> E.Value
encodeOutputForList
    (arg0) =
        E.list (identity)
            [ E.list (encodeTodoItem) arg0
            ]

decodeOutputForList : D.Decoder OutputForList
decodeOutputForList =
        D.map (\a -> (a))
            (D.index 0 (D.map (Maybe.withDefault ([])) (D.maybe (D.list (decodeTodoItem)))))

type alias InputForRetrieve =
    (String)

encodeInputForRetrieve : InputForRetrieve -> E.Value
encodeInputForRetrieve
    (arg0) =
        E.list (identity)
            [ E.string arg0
            ]

decodeInputForRetrieve : D.Decoder InputForRetrieve
decodeInputForRetrieve =
        D.map (\a -> (a))
            (D.index 0 (D.string))

type alias OutputForRetrieve =
    (TodoItem)

encodeOutputForRetrieve : OutputForRetrieve -> E.Value
encodeOutputForRetrieve
    (arg0) =
        E.list (identity)
            [ encodeTodoItem arg0
            ]

decodeOutputForRetrieve : D.Decoder OutputForRetrieve
decodeOutputForRetrieve =
        D.map (\a -> (a))
            (D.index 0 (decodeTodoItem))

type alias InputForUpdate =
    (String, TodoItem)

encodeInputForUpdate : InputForUpdate -> E.Value
encodeInputForUpdate
    (arg0,arg1) =
        E.list (identity)
            [ E.string arg0
            , encodeTodoItem arg1
            ]

decodeInputForUpdate : D.Decoder InputForUpdate
decodeInputForUpdate =
        D.map2 (\arg0 arg1 -> (arg0, arg1))
            (D.index 0 D.string)
            (D.index 1 decodeTodoItem)
    

type alias OutputForUpdate =
    (TodoItem)

encodeOutputForUpdate : OutputForUpdate -> E.Value
encodeOutputForUpdate
    (arg0) =
        E.list (identity)
            [ encodeTodoItem arg0
            ]

decodeOutputForUpdate : D.Decoder OutputForUpdate
decodeOutputForUpdate =
        D.map (\a -> (a))
            (D.index 0 (decodeTodoItem))



callDestroy : Config -> InputForDestroy -> (CallResult OutputForDestroy -> a) -> Cmd a
callDestroy config input mapResult =
    let
        body = Http.jsonBody (encodeInputForDestroy input)
        expect = Http.expectJson (unwrapHttpResult >> mapResult) (decodeCallResult decodeOutputForDestroy)
    in
    Http.request
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/Destroy"
        , body = body
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }

callList : Config -> InputForList -> (CallResult OutputForList -> a) -> Cmd a
callList config input mapResult =
    let
        body = Http.jsonBody (encodeInputForList input)
        expect = Http.expectJson (unwrapHttpResult >> mapResult) (decodeCallResult decodeOutputForList)
    in
    Http.request
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/List"
        , body = body
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }

callRetrieve : Config -> InputForRetrieve -> (CallResult OutputForRetrieve -> a) -> Cmd a
callRetrieve config input mapResult =
    let
        body = Http.jsonBody (encodeInputForRetrieve input)
        expect = Http.expectJson (unwrapHttpResult >> mapResult) (decodeCallResult decodeOutputForRetrieve)
    in
    Http.request
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/Retrieve"
        , body = body
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }

callUpdate : Config -> InputForUpdate -> (CallResult OutputForUpdate -> a) -> Cmd a
callUpdate config input mapResult =
    let
        body = Http.jsonBody (encodeInputForUpdate input)
        expect = Http.expectJson (unwrapHttpResult >> mapResult) (decodeCallResult decodeOutputForUpdate)
    in
    Http.request
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/Update"
        , body = body
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }
