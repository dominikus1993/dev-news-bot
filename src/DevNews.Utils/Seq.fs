namespace DevNews.Utils

module Seq =
    open System.Collections.Generic
    open System.Linq
    
    let paged (pageSize: int)(sequence: _ seq) : _ seq seq =
        seq {
            use enumerator = sequence.GetEnumerator()
            while enumerator.MoveNext() do
                let current = List<_>()
                current.Add(enumerator.Current)
                while (current.Count < pageSize && enumerator.MoveNext()) do
                    current.Add(enumerator.Current)
                yield current.AsEnumerable()
        }
