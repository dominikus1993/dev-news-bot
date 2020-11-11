namespace DevNews.Utils

module Seq =
    let paged (pageSize: int)(sequence: _ seq) =
        seq {
            yield sequence
        }

