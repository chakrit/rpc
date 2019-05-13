module Main exposing (Model, Msg(..), init, main)

import Api exposing (TodoItem)
import Browser exposing (Document)
import Bulma.CDN exposing (stylesheet)
import Bulma.Elements exposing (TablePartition, TitleSize(..), button, buttonModifiers, notification, subtitle, table, tableBody, tableCell, tableCellHead, tableHead, tableModifiers, tableRow, title)
import Bulma.Form exposing (controlButton, controlInput, controlInputModifiers, field)
import Bulma.Layout exposing (SectionSpacing(..), container, section)
import Bulma.Modifiers as Modifiers exposing (Color(..), Size(..), State(..))
import Html exposing (Html, div, li, text, ul)
import Html.Attributes exposing (value)
import Html.Events exposing (onClick, onInput)
import Json.Decode
import RpcUtil as Rpc exposing (CallResult(..), mapResult, translateHttpError)


config : Rpc.Config
config =
    { baseUrl = "http://localhost:7000"
    , headers = []
    }


type alias Flags =
    Json.Decode.Value


type CallState
    = Idle
    | Calling
    | CallFailed String


type alias Model =
    { text : String
    , state : CallState
    , list : List TodoItem
    }


type Msg
    = Input String
    | Create
    | Reset
    | Refresh
    | Delete Int
    | ErrorReply String
    | ListReply (List TodoItem)


onCreateReply =
    mapResult
        { onHttpErr = translateHttpError >> ErrorReply
        , onApiErr = ErrorReply
        , onSuccess = always Reset
        }


onListReply =
    mapResult
        { onHttpErr = translateHttpError >> ErrorReply
        , onApiErr = ErrorReply
        , onSuccess = ListReply
        }


onDeleteReply =
    mapResult
        { onHttpErr = translateHttpError >> ErrorReply
        , onApiErr = ErrorReply
        , onSuccess = always Refresh
        }


init : Flags -> ( Model, Cmd Msg )
init _ =
    ( { text = ""
      , state = Idle
      , list = []
      }
    , Api.callList config () onListReply
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Input str ->
            ( { model | text = str }, Cmd.none )

        Create ->
            ( { model | state = Calling }
            , Api.callCreate config model.text onCreateReply
            )

        Reset ->
            ( { model | state = Calling, text = "" }
            , Api.callList config () onListReply
            )

        Refresh ->
            ( { model | state = Calling }
            , Api.callList config () onListReply
            )

        Delete id ->
            ( { model | state = Calling }
            , Api.callDestroy config id onDeleteReply
            )

        ErrorReply err ->
            ( { model | state = CallFailed err }, Cmd.none )

        ListReply list ->
            ( { model | state = Idle, list = list }, Cmd.none )


view : Model -> Document Msg
view model =
    { title = "RPC example TODO app"
    , body =
        [ stylesheet
        , viewTitle model
        , viewNotification model
        , if List.isEmpty model.list then
            div [] []

          else
            viewList model
        , viewToolbar model
        ]
    }


viewTitle : Model -> Html Msg
viewTitle _ =
    section NotSpaced
        []
        [ container []
            [ title H1 [] [ text "RPC TODO" ]
            ]
        ]


viewNotification : Model -> Html Msg
viewNotification model =
    section NotSpaced
        []
        [ container []
            [ case model.state of
                CallFailed err ->
                    notification Danger [] [ text err ]

                _ ->
                    div [] []
            ]
        ]


viewToolbar : Model -> Html Msg
viewToolbar model =
    let
        input =
            controlInput
                { controlInputModifiers
                    | size = Large
                    , disabled = model.state == Calling
                }
                []
                [ value model.text, onInput Input ]
                []

        button =
            controlButton
                { buttonModifiers
                    | size = Large
                    , color = Modifiers.Success
                    , state =
                        if model.state == Calling then
                            Loading

                        else
                            Blur
                }
                []
                [ onClick Create ]
                [ text "Create" ]
    in
    section NotSpaced
        []
        [ container []
            [ field [] [ input ]
            , button
            ]
        ]


viewList : Model -> Html Msg
viewList model =
    section NotSpaced
        []
        [ container []
            [ subtitle H2
                []
                [ text "Todo List"
                ]
            , table
                { tableModifiers | striped = True, hoverable = True, fullWidth = True }
                []
                [ viewTableHeaders model
                , viewTableBody model
                ]
            ]
        ]


viewTableHeaders : Model -> TablePartition Msg
viewTableHeaders model =
    tableHead []
        [ tableRow False
            []
            [ tableCellHead [] [ text "Item" ]
            , tableCellHead [] [ text "Actions" ]
            ]
        ]


viewTableBody : Model -> Html Msg
viewTableBody model =
    let
        row : TodoItem -> Html Msg
        row todoItem =
            tableRow False
                []
                [ tableCell [] [ text todoItem.description ]
                , tableCell []
                    [ button
                        { buttonModifiers | color = Danger }
                        [ onClick (Delete todoItem.id) ]
                        [ text "Delete" ]
                    ]
                ]
    in
    tableBody [] (List.map row model.list)


main : Program Flags Model Msg
main =
    Browser.document
        { init = init
        , update = update
        , view = view
        , subscriptions = \_ -> Sub.none
        }
