namespace DevNews.Core.Model

type ApplicationError =
    | InsertError

type Article = { Title: string; Link: string }