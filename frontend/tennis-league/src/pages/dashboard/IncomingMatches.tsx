import { useCallback, useEffect, useMemo, useState } from "react";
import { IncomingMatchResponse } from "../../model/dashboard.model";
import { getIncomingMathces } from "../../api/dashboardService";
import { Card } from "primereact/card";
import { formatDate } from "../../helper/date.helper";
import { Button } from "primereact/button";
import { Divider } from "primereact/divider";
import { MatchScoreSidebar } from "../../components/match/MatchScoreSidebar";

interface IncomingMatchesCardProps extends DashboardProps {

}

export default function IncomingMatchesCard({ className = "col-12 md:col-4" }: IncomingMatchesCardProps) {

    const [updateScoreVisible, setUpdateScoreVisible] = useState<boolean>(false);
    const [selectedMatchId, setSelectedMatchId] = useState<string>();
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

    const handleScoreUpdate = (matchId: string) => {
        // Skor güncelleme işlemi burada yapılacak
        setSelectedMatchId(matchId)
        setUpdateScoreVisible(true);
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
                                    onClick={() => handleScoreUpdate(match.matchId)}
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
            <MatchScoreSidebar visible={updateScoreVisible} matchId={selectedMatchId} onHide={() => setUpdateScoreVisible(false)} onSuccess={() => fetchIncomigMatches()} />
        </>


    );
}
