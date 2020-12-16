using System.Collections.Generic;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams.Model
{
    public interface IMicrosoftTeamsMessageContentBody
    {
        string Type { get; }
    }

    public record MicrosoftTeamsAction(string Id, string Title, string Url)
    {
        public string Type { get; } = "Action.OpenUrl";
    }

    public record MicrosoftTeamsTextBlock(string Text) : IMicrosoftTeamsMessageContentBody
    {
        public string Type { get; } = "TextBlock";
    }

    public record MicrosoftTeamsActionSet
        (string Id, IEnumerable<MicrosoftTeamsAction> Actions) : IMicrosoftTeamsMessageContentBody
    {
        public string Type { get; } = "ActionSet";
    }

    public record MicrosoftTeamsMessageContent(IEnumerable<IMicrosoftTeamsMessageContentBody> Body)
    {
        public string Schema { get; } = "http://adaptivecards.io/schemas/adaptive-card.json";
        public string Type { get; } = "AdaptiveCard";
        public string Version { get; } = "1.2";
    }

    public record MicrosoftTeamsAttachment(MicrosoftTeamsMessageContent Content)
    {
        public string ContentType { get; } = "application/vnd.microsoft.card.adaptive";
    }

    public record MicrosoftTeamsMessage(IEnumerable<MicrosoftTeamsAttachment> Attachments)
    {
        public string Type { get; } = "message";

        public MicrosoftTeamsMessage(MicrosoftTeamsAttachment attachment) : this(new[] {attachment})
        {
        }
    }
}