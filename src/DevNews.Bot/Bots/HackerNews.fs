namespace DevNews.Bot

open System.Threading
open DevNews.Core.HackerNews

module HackerNews =
    let rec run (useCase: UseCases.GetNewArticlesAndNotifyUseCase, ct: CancellationToken) = 2

