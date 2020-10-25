using System.Threading.Tasks;
using Akka.Actor;

namespace DevNews.Application.Notifications.Services
{
    public class ActorNotificationService : INotificationService
    {
        private IActorRef _notifierActor;

        public ActorNotificationService(IActorRef notifierActor)
        {
            _notifierActor = notifierActor;
        }

        public async ValueTask Notify()
        {
            throw new System.NotImplementedException();
        }
    }
}