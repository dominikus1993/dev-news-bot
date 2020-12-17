using System.Threading.Tasks;
using DevNews.Infrastructure.Notifications.MicrosoftTeams.Model;
using FluentAssertions;
using Xunit;

namespace DevNews.Infrastructure.Notifications.UnitTests.Teams
{
    public class MessageSerializerUnitTests
    {
        [Fact]
        public async Task TestSerialization()
        {
            var text = new MicrosoftTeamsTextBlock(
                "For Samples and Templates, see https://adaptivecards.io/samples](https://adaptivecards.io/samples)");
            var link = new MicrosoftTeamsAction("5dc0e9c8-0ede-e3d7-1660-b60897c9d2de",
                "adssssssssssssssssssssssssssssssssss", "https://amdesigner.azurewebsites.net");

            var actionSet = new MicrosoftTeamsActionSet("9a6b8098-b3d2-fdb4-f767-8c54f1fc0d1e", new[] {link});
            var content = new MicrosoftTeamsMessageContent(new IMicrosoftTeamsMessageContentBody[] {text, actionSet});
            var att = new MicrosoftTeamsAttachment(content);
            var serializer = new JsonTeamsMessageSerializer();
            var msg = new MicrosoftTeamsMessage(new []{att});
            var subject = await serializer.Serialize(msg);
            subject.Should().Contain("adssssssssssssssssssssssssssssssssss");
        }
    }
}