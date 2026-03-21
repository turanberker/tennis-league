import { Card } from "primereact/card";

interface LeaguePlayersProps {
    leagueId: string;
}

export const LeaguePlayers: React.FC<LeaguePlayersProps> = ({ leagueId }) => {
    return (<Card
        title="Lig Oyuncuları"
    ></Card>);

}