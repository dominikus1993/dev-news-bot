// Learn more about F# at http://docs.microsoft.com/dotnet/fsharp

open System
open Argu

type CliError =
    | ArgumentsNotSpecified of msg: string

type CmdArgs =
    | [<AltCommandLine("-hacker-news")>] Confirm of discordWebHook: int

with
    interface IArgParserTemplate with
        member this.Usage =
            match this with
            | Confirm _ -> "Confirm by workerid"
            | ConfirmAuth -> "Confirm by windows auth"

[<EntryPoint>]
let main argv =
    let message = from "F#" // Call the function
    printfn "Hello world %s" message
    0 // return an integer exit code