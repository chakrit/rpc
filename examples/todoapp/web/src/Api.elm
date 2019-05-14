module Api exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E
import Time
import RpcUtil exposing (Config, CallResult, unwrapHttpResult, decodeCallResult)



type alias TodoItem =
    { ctime : Time.Posix
    , description : String
    , id : Int
    }

encodeTodoItem : TodoItem -> E.Value
encodeTodoItem obj =
    E.object
        [ ( "ctime", (Time.posixToMillis >> toFloat >> (\f -> f/1000.0) >> E.float) obj.ctime )
        , ( "description", E.string obj.description )
        , ( "id", E.int obj.id )
        ]

decodeTodoItem : D.Decoder TodoItem
decodeTodoItem =
    D.map3 TodoItem
            (D.field "ctime" ((D.map ((\f -> f * 1000.0) >> round >> Time.millisToPosix) D.float)))
            (D.field "description" (D.string))
            (D.field "id" (D.int))
    



type alias InputForCreate =
    (String)

encodeInputForCreate : InputForCreate -> E.Value
encodeInputForCreate
    (arg0) =
        E.list (identity)
            [ E.string arg0
            ]

decodeInputForCreate : D.Decoder InputForCreate
decodeInputForCreate =
        D.map (\a -> (a))
            (D.index 0 (D.string))

type alias OutputForCreate =
    (TodoItem)

encodeOutputForCreate : OutputForCreate -> E.Value
encodeOutputForCreate
    (arg0) =
        E.list (identity)
            [ encodeTodoItem arg0
            ]

decodeOutputForCreate : D.Decoder OutputForCreate
decodeOutputForCreate =
        D.map (\a -> (a))
            (D.index 0 (decodeTodoItem))

type alias InputForDestroy =
    (Int)

encodeInputForDestroy : InputForDestroy -> E.Value
encodeInputForDestroy
    (arg0) =
        E.list (identity)
            [ E.int arg0
            ]

decodeInputForDestroy : D.Decoder InputForDestroy
decodeInputForDestroy =
        D.map (\a -> (a))
            (D.index 0 (D.int))

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



callCreate : Config -> InputForCreate -> (CallResult OutputForCreate -> a) -> Cmd a
callCreate config input mapResult =
    let
        body = Http.jsonBody (encodeInputForCreate input)
        expect = Http.expectJson (unwrapHttpResult >> mapResult) (decodeCallResult decodeOutputForCreate)
    in
    Http.request
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/Create"
        , body = body
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }

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
