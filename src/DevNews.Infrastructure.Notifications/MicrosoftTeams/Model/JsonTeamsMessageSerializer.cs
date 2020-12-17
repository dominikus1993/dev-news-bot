using System;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams.Model
{
    public class MicrosoftTeamsMessageContentBodyConverter : JsonConverter<IMicrosoftTeamsMessageContentBody>
    {
        public override IMicrosoftTeamsMessageContentBody? Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            // I need this only for serialization
            throw new NotImplementedException();
        }

        public override void Write(Utf8JsonWriter writer, IMicrosoftTeamsMessageContentBody value, JsonSerializerOptions options)
        {
            switch (value)
            {
                case MicrosoftTeamsActionSet microsoftTeamsActionSet:
                    JsonSerializer.Serialize(writer, microsoftTeamsActionSet, typeof(MicrosoftTeamsActionSet), options);
                    break;
                case MicrosoftTeamsTextBlock microsoftTeamsTextBlock:
                    JsonSerializer.Serialize(writer, microsoftTeamsTextBlock, typeof(MicrosoftTeamsTextBlock), options);
                    break;
                default:
                    throw new ArgumentOutOfRangeException(nameof(value));
            }
        }
    }
    public class JsonTeamsMessageSerializer : ITeamsMessageSerializer
    {
        private static readonly JsonSerializerOptions Options = new()
            {PropertyNamingPolicy = JsonNamingPolicy.CamelCase, Converters = { new MicrosoftTeamsMessageContentBodyConverter() }};
        
        public ValueTask<string> Serialize(MicrosoftTeamsMessage msg)
        {
            return new(JsonSerializer.Serialize(msg, Options));
        }
    }
}