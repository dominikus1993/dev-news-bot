using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using DevNews.Core.Dto;
using DevNews.Core.Model;
using DevNews.Core.Repository;

namespace DevNews.Core.UseCases
{
    public record GetArticlesQuery(int Page, int PageSize);
    
    public class GetArticles
    {
        private readonly IArticlesRepository _articlesRepository;

        public GetArticles(IArticlesRepository articlesRepository)
        {
            _articlesRepository = articlesRepository;
        }

        public async Task<List<ArticleDto>> Execute(GetArticlesQuery query)
        {
            if (query.Page < 0)
            {
                throw new ArgumentOutOfRangeException(nameof(query.Page));
            }

            if (query.PageSize <= 0)
            {
                throw new ArgumentOutOfRangeException(nameof(query.PageSize));
            }
            
            return await _articlesRepository.Get(query.Page, query.PageSize).Select(article => new ArticleDto(article)).ToListAsync();
        }
    }
}