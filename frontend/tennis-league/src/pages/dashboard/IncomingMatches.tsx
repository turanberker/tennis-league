import { useCallback, useEffect, useState } from "react";
import { IncomingMatchResponse } from "../../model/dashboard.model";
import { getIncomingMathces } from "../../api/dashboardService";
import { Card } from "primereact/card";
import { formatDate } from "../../helper/date.helper";
import { Button } from "primereact/button";
import { Divider } from "primereact/divider";
import { MatchScoreSidebar } from "../../components/match/MatchScoreSidebar";
import { MatchScore, MatchSource } from "../../model/match.model";
import { updateFriendlyMatchScore } from "../../api/matchService";
import { updateLeagueMatchScore } from "../../api/leagueService";

interface IncomingMatchesCardProps extends DashboardProps {

}

export default function IncomingMatchesCard({ className = "col-12 md:col-4" }: IncomingMatchesCardProps) {

    const [updateScoreVisible, setUpdateScoreVisible] = useState<boolean>(false);
    const [selectedMatch, setSelectedMatch] = useState<IncomingMatchResponse>();
    const [incomingMatches, setIncomingMatches] = useState<IncomingMatchResponse[]>();
    const [loading, setLoading] = useState(true);

    const fetchIncomigMatches = useCallback(async () => {
        const res = await getIncomingMathces({ limit: 5 })
        setIncomingMatches(res);
        setLoading(false);

    }, []);

    useEffect(() => {
        fetchIncomigMatches();
    }, [fetchIncomigMatches])

    const handleScoreUpdate = (incomingMatch: IncomingMatchResponse) => {
        // Skor güncelleme işlemi burada yapılacak
        setSelectedMatch(incomingMatch)
        setUpdateScoreVisible(true);
    }

    const handleSubmitScore = async (matchId: string, score: MatchScore): Promise<boolean> => {

        let res;
        if (selectedMatch?.source === MatchSource.FRIENDLY) {
            res = await updateFriendlyMatchScore(matchId, score)
        } else if (selectedMatch?.source === MatchSource.LEAGUE) {
            res = await updateLeagueMatchScore(selectedMatch.leagueId!, matchId, score)
        } else {
            throw new Error("Illegal Argument")
        }
        return res != null ? true : false;
    }

    return (
        <>  <div className={className}>
            <Card title="Gelecek Maçlar" style={{ height: '100%' }} className="shadow-2">
                {incomingMatches && incomingMatches.length > 0 ? (
                    incomingMatches.map((match, index) => (
                        <div key={match.matchId}>
                            <div className="flex justify-content-between align-items-center py-2">
                                <div className="flex flex-column gap-1">
                                    {/* Rakip İsmi */}
                                    <span className="font-bold text-900">
                                        {match.oppenentName || "Rakip Bekleniyor"}
                                    </span>

                                    {/* Lig İsmi */}
                                    <span className="text-sm text-500">
                                        <i className="pi pi-trophy mr-1 text-xs"></i>
                                        {match.leagueName}
                                    </span>

                                    {/* Tarih */}
                                    <span className="text-xs text-600">
                                        <i className="pi pi-calendar mr-1 text-xs"></i>
                                        {formatDate(match.matchDate)}
                                    </span>
                                </div>

                                {/* Skor Gir Butonu */}
                                <Button
                                    icon="pi pi-pencil"
                                    className="p-button-rounded p-button-text p-button-sm"
                                    tooltip="Skor Gir"
                                    onClick={() => handleScoreUpdate(match)}
                                />
                            </div>
                            {/* Son eleman değilse araya çizgi koy */}
                            {index !== incomingMatches.length - 1 && <Divider className="my-2" />}
                        </div>
                    ))
                ) : (
                    <div className="text-center py-4 text-500">
                        {loading ? "Yükleniyor..." : "Yakın zamanda maçınız bulunmuyor."}
                    </div>
                )}
            </Card>
        </div>
            <MatchScoreSidebar visible={updateScoreVisible} matchId={selectedMatch?.matchId} onHide={() => setUpdateScoreVisible(false)}
                submitMatchScore={handleSubmitScore} onSuccess={() => fetchIncomigMatches()} />
        </>


    );
}
