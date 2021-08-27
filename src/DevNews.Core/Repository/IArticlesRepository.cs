using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using DevNews.Core.Model;
using LanguageExt;

namespace DevNews.Core.Repository
{
    public interface IArticlesRepository
    {
        Task<bool> Exists(Article article, CancellationToken cancellationToken = default);
        Task<Either<Exception, Unit>> InsertMany(IEnumerable<Article> articles, CancellationToken cancellationToken = default);
        IAsyncEnumerable<Article> Get(int page, int pageSize, CancellationToken cancellationToken = default);
        Task<long> Count(CancellationToken cancellationToken = default);
    }
}