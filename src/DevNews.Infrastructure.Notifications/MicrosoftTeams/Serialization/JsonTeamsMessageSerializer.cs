using DevNews.Infrastructure.Notifications.MicrosoftTeams.Model;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams.Serialization
{
    internal class JsonTeamsMessageSerializer : ITeamsMessageSerializer
    {
        private static readonly JsonSerializerOptions Options = new()
        { IgnoreNullValues = true, PropertyNamingPolicy = JsonNamingPolicy.CamelCase };

        public ValueTask<string> Serialize(MicrosoftTeamsMessage msg)
        {
            return new(JsonSerializer.Serialize(msg, Options));
        }
    }
}
