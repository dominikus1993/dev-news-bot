using System;
using System.Threading.Tasks;
using DevNews.Core.Repository;

namespace DevNews.Core.UseCases
{
    public class GetArticlesPagesQuantity
    {
        private IArticlesRepository _articlesRepository;

        public GetArticlesPagesQuantity(IArticlesRepository articlesRepository)
        {
            _articlesRepository = articlesRepository;
        }

        public async Task<long> Execute(int pageSize)
        {
            if (pageSize <= 0)
            {
                throw new ArgumentOutOfRangeException(nameof(pageSize), "Page Size should be greater or equal 1");
            }

            var count = await _articlesRepository.Count();
            if (count == 0)
            {
                return 0;
            }

            return (long) Math.Ceiling((double) count / (double) pageSize);
        }
    }
}