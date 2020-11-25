// Learn more about F# at http://docs.microsoft.com/dotnet/fsharp

open System
open Argu

type CliError =
    | ArgumentsNotSpecified of msg: string

type CmdArgs =
    | [<AltCommandLine("-parse")>] Parse of discordWebHookUrl: string * mongoConnection:string

with
    interface IArgParserTemplate with
        member this.Usage =
            match this with
            | Parse _ -> "Parse articles and notify users"

[<EntryPoint>]
let main argv =
    0 // return an integer exit code