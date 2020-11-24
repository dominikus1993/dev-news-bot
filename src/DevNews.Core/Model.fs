namespace DevNews.Core.Model

type ApplicationError =
    | InsertArticleError
    | NotifyUserError

type Article = { Title: string; Link: string; Source: string }