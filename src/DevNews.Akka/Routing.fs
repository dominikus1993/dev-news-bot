namespace DevNews.Shared.Routing

type ActorMetaData  = { Name: string; Parent: ActorMetaData option; Path: string }


module ActorMetaData  =
    let Create(name, parent) =
        let parentPath = match parent with
                         | Some p -> p.Path
                         | None -> "/user"
        { Name = name; Parent = parent; Path = (sprintf "%s/%s" parentPath name) }
        
    let CreateTopLevel(name) =
        Create(name, None)