module Api exposing (..)

-- <auto-generated />
-- @generated by github.com/chakrit/rpc

import Http
import Json.Decode as D
import Json.Encode as E
import Dict exposing (Dict)
import Task exposing (Task)
import Time exposing (Posix)
import Bytes exposing (Bytes)
import Bytes.Encode
import RpcUtil exposing (Config, RpcError, RpcResult, decodeApply, fromHttpResult)



type alias TodoItem =
    { ctime : Posix
    , description : String
    , id : Int
    , metadata : String
    , state : State
    }

defaultTodoItem : TodoItem
defaultTodoItem =
    { ctime = Time.millisToPosix 0
    , description = ""
    , id = 0
    , metadata = ""
    , state = defaultState
    }

encodeTodoItem : TodoItem -> E.Value
encodeTodoItem obj =
    E.object
        [ ( "ctime", (Time.posixToMillis >> toFloat >> (\f -> f/1000.0) >> E.float) obj.ctime )
        , ( "description", E.string obj.description )
        , ( "id", E.int obj.id )
        , ( "metadata", E.string obj.metadata )
        , ( "state", encodeState obj.state )
        ]

decodeTodoItem : D.Decoder TodoItem
decodeTodoItem =
    D.map5 TodoItem
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
                (D.string
                    |> D.field "metadata"
                    |> D.maybe
                    |> D.map (Maybe.withDefault (""))
                )
                (decodeState
                    |> D.field "state"
                    |> D.maybe
                    |> D.map (Maybe.withDefault (defaultState))
                )
    



type State
    = New
    | InProgress
    | Overdue
    | Completed

allState : List State
allState =
    [ New
    , InProgress
    , Overdue
    , Completed
    ]

pairsOfState : List ( String, State )
pairsOfState =
    [ ( "new", New )
    , ( "in-progress", InProgress )
    , ( "overdue", Overdue )
    , ( "completed", Completed )
    ]

titlePairsOfState : List ( String, String )
titlePairsOfState =
    [ ( "new", "New" )
    , ( "in-progress", "In Progress" )
    , ( "overdue", "Overdue" )
    , ( "completed", "Completed" )
    ]

stringToState : String -> Maybe State
stringToState str =
    case str of
        "new" ->
            Just New
        "in-progress" ->
            Just InProgress
        "overdue" ->
            Just Overdue
        "completed" ->
            Just Completed
        _ ->
            Nothing

stringFromState : State -> String
stringFromState v =
    case v of
        New ->
            "new"
        InProgress ->
            "in-progress"
        Overdue ->
            "overdue"
        Completed ->
            "completed"

titleStringFromState : State -> String
titleStringFromState v =
    case v of
        New ->
            "New"
        InProgress ->
            "In Progress"
        Overdue ->
            "Overdue"
        Completed ->
            "Completed"

defaultState =
    New

encodeState : State -> E.Value
encodeState =
    stringFromState >> E.string

decodeState : D.Decoder State
decodeState =
    D.string
        |> D.map stringToState
        |> D.map (Maybe.withDefault defaultState)



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

type alias InputForUpdateState =
    (Int, State)

encodeInputForUpdateState : InputForUpdateState -> E.Value
encodeInputForUpdateState
    (arg0,arg1) =
        E.list (identity)
            [ E.int arg0
            , encodeState arg1
            ]

decodeInputForUpdateState : D.Decoder InputForUpdateState
decodeInputForUpdateState =
        D.map2 (\arg0 arg1 -> (arg0, arg1))
            (D.int
                |> D.index 0
                |> D.maybe
                |> D.map (Maybe.withDefault (0))
            )
            (decodeState
                |> D.index 1
                |> D.maybe
                |> D.map (Maybe.withDefault (defaultState))
            )
    

type alias OutputForUpdateState =
    (TodoItem)

encodeOutputForUpdateState : OutputForUpdateState -> E.Value
encodeOutputForUpdateState
    (arg0) =
        E.list (identity)
            [ encodeTodoItem arg0
            ]

decodeOutputForUpdateState : D.Decoder OutputForUpdateState
decodeOutputForUpdateState =
        decodeTodoItem
            |> D.index 0
            |> D.maybe
            |> D.map (Maybe.withDefault (defaultTodoItem))
            |> D.map (\a -> (a))



callCreateTask : Config -> InputForCreate -> Task RpcError OutputForCreate
callCreateTask config input =
    let
        body =
            Http.jsonBody (encodeInputForCreate input)

        resolver =
            RpcUtil.resolver decodeOutputForCreate
    in
    Http.task
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/Create"
        , body = body
        , resolver = resolver
        , timeout = Nothing
        }


callCreate : Config -> InputForCreate -> (RpcResult OutputForCreate -> a) -> Cmd a
callCreate config input mapResult =
    let
        body = Http.jsonBody (encodeInputForCreate input)
        expect = Http.expectJson (fromHttpResult >> mapResult) (RpcUtil.decoder decodeOutputForCreate)
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

callDestroyTask : Config -> InputForDestroy -> Task RpcError OutputForDestroy
callDestroyTask config input =
    let
        body =
            Http.jsonBody (encodeInputForDestroy input)

        resolver =
            RpcUtil.resolver decodeOutputForDestroy
    in
    Http.task
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/Destroy"
        , body = body
        , resolver = resolver
        , timeout = Nothing
        }


callDestroy : Config -> InputForDestroy -> (RpcResult OutputForDestroy -> a) -> Cmd a
callDestroy config input mapResult =
    let
        body = Http.jsonBody (encodeInputForDestroy input)
        expect = Http.expectJson (fromHttpResult >> mapResult) (RpcUtil.decoder decodeOutputForDestroy)
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

callListTask : Config -> InputForList -> Task RpcError OutputForList
callListTask config input =
    let
        body =
            Http.jsonBody (encodeInputForList input)

        resolver =
            RpcUtil.resolver decodeOutputForList
    in
    Http.task
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/List"
        , body = body
        , resolver = resolver
        , timeout = Nothing
        }


callList : Config -> InputForList -> (RpcResult OutputForList -> a) -> Cmd a
callList config input mapResult =
    let
        body = Http.jsonBody (encodeInputForList input)
        expect = Http.expectJson (fromHttpResult >> mapResult) (RpcUtil.decoder decodeOutputForList)
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

callUpdateStateTask : Config -> InputForUpdateState -> Task RpcError OutputForUpdateState
callUpdateStateTask config input =
    let
        body =
            Http.jsonBody (encodeInputForUpdateState input)

        resolver =
            RpcUtil.resolver decodeOutputForUpdateState
    in
    Http.task
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/UpdateState"
        , body = body
        , resolver = resolver
        , timeout = Nothing
        }


callUpdateState : Config -> InputForUpdateState -> (RpcResult OutputForUpdateState -> a) -> Cmd a
callUpdateState config input mapResult =
    let
        body = Http.jsonBody (encodeInputForUpdateState input)
        expect = Http.expectJson (fromHttpResult >> mapResult) (RpcUtil.decoder decodeOutputForUpdateState)
    in
    Http.request
        { method = "POST"
        , headers = config.headers
        , url = config.baseUrl ++ "/api/UpdateState"
        , body = body
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }

