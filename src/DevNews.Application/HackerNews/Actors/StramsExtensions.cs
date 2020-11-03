using System.Collections.Generic;
using Akka.Streams;
using Akka.Streams.Dsl;
using DevNews.Core.Model;

namespace DevNews.Application.HackerNews.Actors
{
    public static class StreamsExtensions
    {
        public static dynamic A()
        {
            var a = GraphDsl.Create(builder =>
            {
                var broadcast = builder.Add(new Broadcast<>(2));
                    
                return FlowShape<List<Article>>
            })
        }
    }
}