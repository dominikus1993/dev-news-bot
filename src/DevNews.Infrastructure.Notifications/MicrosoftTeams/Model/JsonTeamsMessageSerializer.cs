using System.Text.Json;
using System.Threading.Tasks;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams.Model
{
    public class JsonTeamsMessageSerializer : ITeamsMessageSerializer
    {
        private static readonly JsonSerializerOptions Options = new()
            {IgnoreNullValues = true, PropertyNamingPolicy = JsonNamingPolicy.CamelCase};
        
        public ValueTask<string> Serialize(MicrosoftTeamsMessage msg)
        {
            return new(JsonSerializer.Serialize(msg, Options));
        }
    }
}