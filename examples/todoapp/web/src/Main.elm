module Main exposing (Model, Msg(..), init, main)

import Api exposing (TodoItem)
import Browser exposing (Document)
import Bulma.CDN exposing (stylesheet)
import Bulma.Elements exposing (TablePartition, TitleSize(..), button, buttonModifiers, notification, subtitle, table, tableBody, tableCell, tableCellHead, tableHead, tableModifiers, tableRow, title)
import Bulma.Form exposing (controlButton, controlInput, controlInputModifiers, field)
import Bulma.Layout exposing (SectionSpacing(..), container, section)
import Bulma.Modifiers as Modifiers exposing (Color(..), Size(..), State(..))
import Html exposing (Html, div, text)
import Html.Attributes exposing (value)
import Html.Events exposing (onClick, onInput)
import Json.Decode
import RpcUtil as Rpc exposing (CallResult(..), b64StringFromBytes, mapResult, translateHttpError)
import Task
import Time exposing (Month(..), utc)


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
    { zone : Time.Zone
    , text : String
    , state : CallState
    , list : List TodoItem
    }


type Msg
    = TimeZone Time.Zone
    | Input String
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
    ( { zone = utc
      , text = ""
      , state = Idle
      , list = []
      }
    , Cmd.batch
        [ Api.callList config () onListReply
        , Time.here |> Task.perform TimeZone
        ]
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        TimeZone zone ->
            ( { model | zone = zone }, Cmd.none )

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
            , tableCellHead [] [ text "Created" ]
            , tableCellHead [] [ text "Actions" ]
            , tableCellHead [] [ text "Metadata" ]
            ]
        ]


viewTableBody : Model -> Html Msg
viewTableBody model =
    let
        rowTime : TodoItem -> Html Msg
        rowTime todoItem =
            let
                t =
                    todoItem.ctime

                ( day, month, year ) =
                    ( Time.toDay model.zone t |> String.fromInt
                    , Time.toMonth model.zone t |> monthName
                    , Time.toYear model.zone t |> String.fromInt
                    )

                ( hour, minute ) =
                    ( Time.toHour model.zone t |> String.fromInt
                    , Time.toMinute model.zone t |> String.fromInt
                    )
            in
            text
                (day
                    ++ " "
                    ++ month
                    ++ " "
                    ++ year
                    ++ " "
                    ++ hour
                    ++ ":"
                    ++ minute
                )

        row : TodoItem -> Html Msg
        row todoItem =
            tableRow False
                []
                [ tableCell [] [ text todoItem.description ]
                , tableCell [] [ rowTime todoItem ]
                , tableCell []
                    [ button
                        { buttonModifiers | color = Danger }
                        [ onClick (Delete todoItem.id) ]
                        [ text "Delete" ]
                    ]
                , tableCell [] [ text (b64StringFromBytes todoItem.metadata |> Maybe.withDefault "(n/a)") ]
                ]
    in
    tableBody [] (List.map row model.list)


monthName : Time.Month -> String
monthName m =
    case m of
        Jan ->
            "Jan"

        Feb ->
            "Feb"

        Mar ->
            "Mar"

        Apr ->
            "Apr"

        May ->
            "May"

        Jun ->
            "Jun"

        Jul ->
            "Jul"

        Aug ->
            "Aug"

        Sep ->
            "Sep"

        Oct ->
            "Oct"

        Nov ->
            "Nov"

        Dec ->
            "Dec"


main : Program Flags Model Msg
main =
    Browser.document
        { init = init
        , update = update
        , view = view
        , subscriptions = \_ -> Sub.none
        }
