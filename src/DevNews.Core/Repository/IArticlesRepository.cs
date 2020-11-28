using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using DevNews.Core.Model;
using LanguageExt;

namespace DevNews.Core.Repository
{
    public interface IArticlesRepository
    {
        Task<bool> Exists(Article article);
        Task<Either<Exception, Unit>> InsertMany(IEnumerable<Article> articles);
    }
}