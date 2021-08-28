using DevNews.Core.Model;
using System;
using System.Collections.Generic;
using System.Linq;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams.Model
{
    internal interface IMicrosoftTeamsMessageContentBody
    {
        string Type { get; }
    }

    internal record MicrosoftTeamsAction(string Id, string Title, string Url)
    {
        public string Type { get; } = "Action.OpenUrl";
    }

    internal record MicrosoftTeamsTextBlock(string Text) : IMicrosoftTeamsMessageContentBody
    {
        public string Type { get; } = "TextBlock";
    }

    internal record MicrosoftTeamsActionSet
        (string Id, IEnumerable<MicrosoftTeamsAction> Actions) : IMicrosoftTeamsMessageContentBody
    {
        public string Type { get; } = "ActionSet";
    }

    internal record MicrosoftTeamsMessageContent(IEnumerable<IMicrosoftTeamsMessageContentBody> Body)
    {
        public string Schema { get; } = "http://adaptivecards.io/schemas/adaptive-card.json";
        public string Type { get; } = "AdaptiveCard";
        public string Version { get; } = "1.2";
    }

    internal record MicrosoftTeamsAttachment(MicrosoftTeamsMessageContent Content)
    {
        public string ContentType { get; } = "application/vnd.microsoft.card.adaptive";
    }

    internal record MicrosoftTeamsMessage(IEnumerable<MicrosoftTeamsAttachment> Attachments)
    {
        public string Type { get; } = "message";

        public MicrosoftTeamsMessage(MicrosoftTeamsAttachment attachment) : this(new[] { attachment })
        {
        }

        public static MicrosoftTeamsMessage From(IEnumerable<Article> articles)
        {
            var text = new MicrosoftTeamsTextBlock("Witam serdecznie, oto nowe newsy");
            var links = articles.Select(article =>
                new MicrosoftTeamsAction(Guid.NewGuid().ToString(), article.Title, article.Link)).ToList();
            var action = new MicrosoftTeamsActionSet(Guid.NewGuid().ToString(), links);
            var content = new MicrosoftTeamsMessageContent(new IMicrosoftTeamsMessageContentBody[] { text, action });
            var att = new MicrosoftTeamsAttachment(content);
            return new MicrosoftTeamsMessage(att);
        }
    }
}
