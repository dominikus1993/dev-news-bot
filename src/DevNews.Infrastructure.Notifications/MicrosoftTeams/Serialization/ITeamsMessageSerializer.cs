using DevNews.Infrastructure.Notifications.MicrosoftTeams.Model;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams.Serialization
{
    internal interface ITeamsMessageSerializer
    {
        ValueTask<string> Serialize(MicrosoftTeamsMessage msg);
    }
}
