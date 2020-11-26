// Learn more about F# at http://docs.microsoft.com/dotnet/fsharp

open System
open Argu
open DevNews.Core.Model
open DevNews.Cli

type CliError =
    | ArgumentsNotSpecified of msg: string

type CmdArgs =
    | [<AltCommandLine("-parse")>] Parse of mongoConnection:string * discordWebHookUrl: string 

with
    interface IArgParserTemplate with
        member this.Usage =
            match this with
            | Parse _ -> "Parse articles and notify users"


let parse (mongoConnection:string, discordWebHookUrl: string) =
    printfn $"{mongoConnection} {discordWebHookUrl}"
    let useCase = Articles.CompositionRoot.create (mongoConnection, discordWebHookUrl)
    async { 
        do! useCase.ParseArticlesAndNotify({ ArticlesQuantity = 5 })
        return Ok(());
    }

let getExitCode result =
    async {
        match! result with
        | Ok (res) -> 
            printfn "%A" res
            return 0
        | Error err ->
            match err with
            | ArgumentsNotSpecified(err) -> 
                printfn "%s" (err)
                return 1
    }

[<EntryPoint>]
let main argv =
    let errorHandler = ProcessExiter(colorizer = function ErrorCode.HelpText -> None | _ -> Some ConsoleColor.Red)
    let parser = ArgumentParser.Create<CmdArgs>(programName = "dev-news-cli", errorHandler = errorHandler)
    let res = match parser.ParseCommandLine argv with
              | p when p.Contains(Parse) ->  parse (p.GetResult(Parse))
              | _ ->
                  async { return Error(ArgumentsNotSpecified(parser.PrintUsage()))  }
              |> getExitCode
    res |> Async.RunSynchronously