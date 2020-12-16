using System.Threading.Tasks;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams.Model
{
    public interface ITeamsMessageSerializer
    {
        ValueTask<string> Serialize(MicrosoftTeamsMessage msg);
    }
}