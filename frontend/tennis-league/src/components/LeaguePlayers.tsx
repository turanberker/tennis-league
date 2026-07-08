import { Card } from "primereact/card";
import {getPlayers, getTeams} from "../api/mail/leagueService";
import {useEffect, useState} from "react";
import {useLeague} from "../hooks/useLeague";
import {} from "../model/team.model";
import {SingleLeagueAttendanceResponse} from "../model/league.model";
import {DataTable} from "primereact/datatable";
import {Column} from "primereact/column";
import Guard from "../helper/Guard";
import {Role} from "../model/user.model";
import {Button} from "primereact/button";

interface LeaguePlayersProps {
    leagueId: string;
}

export const LeaguePlayers: React.FC<LeaguePlayersProps> = ({ leagueId }) => {
    const { data: league, updateLeagueCache } = useLeague(leagueId)
    const [players, setPlayers] = useState<SingleLeagueAttendanceResponse[]>([]);
    const [loading, setLoading] = useState<boolean>(false);

    const loadTeams = async (): Promise<void> => {
        if (!leagueId) return;

        setLoading(true);
        const res = await getPlayers(leagueId);
        setPlayers(res);
        setLoading(false);
    };

    useEffect(() => {
        loadTeams();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [leagueId]);

    const header = () => {
        return (
            <Guard allowedRoles={[Role.ADMIN, Role.COORDINATOR]}>
                <div className="flex justify-content-end">
                    <Button
                        label="Yeni Takım"
                        icon="pi pi-plus"
                        size="small"
                        onClick={() => {
                           // loadPlayers();
                           // setCreateDialogVisible(true)
                        }}
                    />
                </div>
            </Guard>
        );
    };

    return (<Card
        title="Oyuncular"
    >
        <DataTable
            value={players}
            loading={loading}
            emptyMessage="Oyuncu bulunamadı"
            tableStyle={{ minWidth: '50rem' }}
            header={header}
            key="id"
        >
            <Column body="firstname" header="Adı" />
            <Column body="surname" header="Soyadı" />
            <Column field="power" header="Gücü" />
        </DataTable>

    </Card>);

}