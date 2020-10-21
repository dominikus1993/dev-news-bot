using System.Threading.Tasks;

namespace DevNews.WebHooks.Application.Services
{
    public interface IWebHookNotifier
    {
        Task Notify();
    }
}