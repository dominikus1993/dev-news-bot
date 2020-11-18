namespace DevNews.Utils

module AsyncResult =
    let bind (f : 'a -> Async<Result<'b, 'c>>) (a : Async<Result<'a, 'c>>)  : Async<Result<'b, 'c>> = async {
        match! a with
        | Ok value ->
            let next : Async<Result<'b, 'c>> = f value
            return! next
        | Error err -> return (Error err)
    }

    let map (f : 'a -> Async<'b>) (a : Async<Result<'a, 'c>>)  : Async<Result<'b, 'c>> = async {
        match! a with
        | Ok value ->
            let! next = f value
            return Ok(next)
        | Error err -> return (Error err)
    }    

module AsyncOption =
    
    let bind (f : 'a -> Async<Option<'b>>) (a : Async<Option<'a>>)  : Async<Option<'b>> = async {
        match! a with
        | Some value ->
            let next : Async<Option<'b>> = f value
            return! next
        | None -> return None
    }

    let ifSome (f : 'a -> Async<unit>) (a : Async<Option<'a>>)  : Async<unit> = async {
        match! a with
        | Some value ->
            do! f value
            return ()
        | None -> return ()
    }