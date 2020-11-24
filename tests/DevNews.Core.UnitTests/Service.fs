namespace DevNews.Core.UnitTests

open FSharp.Control


module ServiceTests =
    open Xunit
    open DevNews.Core.Model
    open DevNews.Core
    open FsUnit.Xunit
    
    let private fakeArticleParser() =
        asyncSeq {
            yield { Link = "a"; Title = "a"; Source = "a"}
            yield { Link = "b"; Title = "b"; Source = "b"}
            yield { Link = "c"; Title = "c"; Source = "c"}
        }
    
    let fakeExistenceChecker(article: Article) =
        async {
            return article.Title = "a"
        }
    
    [<Fact>]
    let ``Test sequence`` () =
        async {
            let provider = Service.provideNewArticles(fakeArticleParser)(fakeExistenceChecker)
            let! subject =  provider() |> AsyncSeq.toListAsync
            let expected = [{ Link = "b"; Title = "b"; Source = "b"}; { Link = "c"; Title = "c"; Source = "c"};]
            subject |> should equal  expected
        }

