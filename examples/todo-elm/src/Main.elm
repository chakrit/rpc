module Main exposing (Flags, Model, Msg(..), main)

import Api
import Browser
import Html exposing (Html, button, div, input, li, text, ul)
import Html.Attributes exposing (value)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Encode
import RpcUtil exposing (CallResult(..))


type alias Flags =
    Json.Encode.Value


type Msg
    = NoOp
    | SetError String
    | SetDirty
    | RefreshItems (List Api.TodoItem)
    | InputNewItem String
    | SaveNewItem
    | DeleteItem String


type alias Model =
    { todoItems : List Api.TodoItem
    , newItem : Api.TodoItem
    , dirty : Bool
    , error : String
    }


init : Flags -> ( Model, Cmd Msg )
init flags =
    ( { todoItems = []
      , newItem =
            { description = "(new item)"
            , id = ""
            , done = False
            }
      , dirty = False
      , error = ""
      }
    , Cmd.none
    )


view : Model -> Browser.Document Msg
view model =
    { title = "Elm TODO w/ RPC"
    , body =
        [ viewErrorBar model
        , viewNewItemBar model
        , viewItemList model
        ]
    }


viewErrorBar : Model -> Html Msg
viewErrorBar model =
    div []
        (if String.isEmpty model.error then
            []

         else
            [ text model.error ]
        )


viewNewItemBar : Model -> Html Msg
viewNewItemBar model =
    div []
        [ input [ onInput InputNewItem, value model.newItem.description ] []
        , button [ onClick SaveNewItem ] [ text "Save" ]
        ]


viewItemList : Model -> Html Msg
viewItemList model =
    let
        listItem : Api.TodoItem -> Html Msg
        listItem todoItem =
            li []
                [ text todoItem.description
                , button [ onClick (DeleteItem todoItem.id) ] [ text "Delete" ]
                ]
    in
    div [] [ ul [] (List.map listItem model.todoItems) ]


rpcConfig : RpcUtil.Config
rpcConfig =
    { baseUrl = "http://0.0.0.0:9999"
    , headers = []
    }


translateHttpError : Http.Error -> String
translateHttpError err =
    case err of
        Http.BadUrl str ->
            "Invalid URL: " ++ str

        Http.Timeout ->
            "Timeout"

        Http.NetworkError ->
            "Network Error"

        Http.BadStatus code ->
            "Bad Status: " ++ String.fromInt code

        Http.BadBody body ->
            "Bad Body: " ++ body


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp ->
            ( model, Cmd.none )

        SetError str ->
            ( { model | error = str }, Cmd.none )

        SetDirty ->
            let
                response : RpcUtil.CallResult (List Api.TodoItem) -> Msg
                response result =
                    case result of
                        HttpError httpErr ->
                            SetError ("HTTP" ++ translateHttpError httpErr)

                        ApiError errStr ->
                            SetError errStr

                        Success items ->
                            RefreshItems items
            in
            ( model
            , Api.callList rpcConfig () response
            )

        RefreshItems list ->
            ( { model | todoItems = list }, Cmd.none )

        InputNewItem str ->
            let
                newItem =
                    model.newItem
            in
            ( { model | newItem = { id = "", description = str, done = False } }, Cmd.none )

        SaveNewItem ->
            let
                response : CallResult Api.TodoItem -> Msg
                response result =
                    case result of
                        HttpError httpErr ->
                            SetError (translateHttpError httpErr)

                        ApiError errStr ->
                            SetError errStr

                        Success _ ->
                            SetDirty
            in
            ( { model | newItem = { id = "", description = "", done = False } }
            , Api.callUpdate rpcConfig ( model.newItem.description, model.newItem ) response
            )

        DeleteItem itemId ->
            let
                response : CallResult Api.TodoItem -> Msg
                response result =
                    case result of
                        HttpError httpErr ->
                            SetError (translateHttpError httpErr)

                        ApiError errStr ->
                            SetError errStr

                        Success _ ->
                            SetDirty
            in
            ( model, Api.callDestroy rpcConfig itemId response )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


main : Program Flags Model Msg
main =
    Browser.document
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
