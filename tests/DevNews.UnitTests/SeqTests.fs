module DevNews.UnitTests.SeqTests
open Xunit
open DevNews.Utils
open FsUnit.Xunit

[<Fact>]
let ``Test empty sequence`` () =
    let emptySeq: _ seq = Seq.empty
    let subject = emptySeq |> Seq.paged 10
    subject |> should be Empty
    
    
[<Fact>]
let ``Test sequence`` () =
    let emptySeq: _ seq = seq [|1;2;3;4;5;|]
    let subject = emptySeq |> Seq.paged 2 |> Seq.map(Seq.toArray) |> Seq.toArray
    let expected = [|[|1;2|]; [|3;4|]; [|5|]|]
    subject |> should equal  expected