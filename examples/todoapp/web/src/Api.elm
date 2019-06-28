module Api exposing (..)

import Http
import Json.Decode as D
import Json.Encode as E
import Dict exposing (Dict)
import Time exposing (Posix)
import Bytes exposing (Bytes)
import Bytes.Encode
import RpcUtil exposing (Config, CallResult, unwrapHttpResult, decodeCallResult, decodeApply)



type alias TodoItem =
    { ctime : Posix
    , description : String
    , id : Int
    , metadata : Bytes
    }

defaultTodoItem : TodoItem
defaultTodoItem =
    { ctime = Time.millisToPosix 0
    , description = ""
    , id = 0
    , metadata = RpcUtil.emptyBytes
    }

encodeTodoItem : TodoItem -> E.Value
encodeTodoItem obj =
    E.object
        [ ( "ctime", (Time.posixToMillis >> toFloat >> (\f -> f/1000.0) >> E.float) obj.ctime )
        , ( "description", E.string obj.description )
        , ( "id", E.int obj.id )
        , ( "metadata", (RpcUtil.b64StringFromBytes >> Maybe.withDefault "" >> E.string) obj.metadata )
        ]

decodeTodoItem : D.Decoder TodoItem
decodeTodoItem =
    D.map4 TodoItem
                ((D.map ((\f -> f * 1000.0) >> round >> Time.millisToPosix) D.float)
                    |> D.field "ctime"
                    |> D.maybe
                    |> D.map (Maybe.withDefault (Time.millisToPosix 0))
                )
                (D.string
                    |> D.field "description"
                    |> D.maybe
                    |> D.map (Maybe.withDefault (""))
                )
                (D.int
                    |> D.field "id"
                    |> D.maybe
                    |> D.map (Maybe.withDefault (0))
                )
                ((D.map (Maybe.withDefault "" >> RpcUtil.b64StringToBytes >> Maybe.withDefault (Bytes.Encode.encode (Bytes.Encode.string ""))) (D.maybe D.string))
                    |> D.field "metadata"
                    |> D.maybe
                    |> D.map (Maybe.withDefault (RpcUtil.emptyBytes))
                )
    



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
        D.string
            |> D.index 0
            |> D.maybe
            |> D.map (Maybe.withDefault (""))
            |> D.map (\a -> (a))

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
        decodeTodoItem
            |> D.index 0
            |> D.maybe
            |> D.map (Maybe.withDefault (defaultTodoItem))
            |> D.map (\a -> (a))

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
        D.int
            |> D.index 0
            |> D.maybe
            |> D.map (Maybe.withDefault (0))
            |> D.map (\a -> (a))

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
        decodeTodoItem
            |> D.index 0
            |> D.maybe
            |> D.map (Maybe.withDefault (defaultTodoItem))
            |> D.map (\a -> (a))

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
        D.list (decodeTodoItem)
            |> D.index 0
            |> D.maybe
            |> D.map (Maybe.withDefault ([]))
            |> D.map (\a -> (a))



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

