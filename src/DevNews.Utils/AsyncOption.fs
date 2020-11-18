namespace DevNews.Utils

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