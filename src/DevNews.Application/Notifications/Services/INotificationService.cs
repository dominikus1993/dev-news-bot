using System.Threading.Tasks;

namespace DevNews.Application.Notifications.Services
{
    public interface INotificationService
    {
        ValueTask Notify();
    }
}