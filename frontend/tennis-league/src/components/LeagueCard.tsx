import React, { } from 'react';
import { League_Category_Labels, League_Format_Labels, League_Process_Type_Labels, League_Status_Labels } from '../model/league.model';
import { Skeleton } from 'primereact/skeleton';
import { Card } from 'primereact/card';
import { formatDate } from '../helper/date.helper';
import { useLeague } from '../hooks/useLeague';


interface LeagueCardProps {
    id: string;
}

export const LeagueCard: React.FC<LeagueCardProps> = ({ id }) => {

    const { data: league, isLoading } = useLeague(id);


    if (isLoading) {
        return <Skeleton width="100%" height="150px" />;
    }

    if (!league) {
        return <div>Lig bilgisi bulunamadı.</div>;
    }

    return (
        <Card className="mb-2 shadow-2" content='p-0' title={league.name} pt={{
            body: { className: 'p-3' },      // Body padding'ini daralttık
            content: { className: 'p-0' }    // İçerik padding'ini tamamen sıfırladık
        }}>
            <div className="grid">
                {/* 1. KOLON: Temel Bilgiler */}
                <div className="col-12 md:col-4 border-right-1 border-200">
                    <div className="flex flex-column gap-3 p-2">
                        <div className="flex justify-content-between align-items-center">
                            <span className="font-bold">Kategori:</span>
                            <span >{League_Category_Labels[league.category]}</span>
                        </div>
                        <div className="flex justify-content-between align-items-center">
                            <span className="font-bold">Format:</span>
                            <span >{League_Format_Labels[league.format]}</span>
                        </div>
                    </div>
                </div>

                {/* 2. KOLON: Durum ve Katılım */}
                <div className="col-12 md:col-4 border-right-1 border-200">
                    <div className="flex flex-column gap-3 p-2">
                        <div className="flex justify-content-between align-items-center">
                            <span className="font-bold">İşleyiş:</span>
                            <span >{League_Process_Type_Labels[league.processType]}</span>
                        </div>
                        <div className="flex justify-content-between align-items-center">
                            <span className="font-bold">Katılımcı:</span>
                            <span >{league.totalAttentance}</span>
                        </div>
                        <div className="flex justify-content-between align-items-center">
                            <span className="font-bold">Durum:</span>
                            <span >{League_Status_Labels[league.status]}</span>

                        </div>
                    </div>
                </div>

                {/* 3. KOLON: Tarihler */}
                <div className="col-12 md:col-4">
                    <div className="flex flex-column gap-3 p-2">
                        <div className="flex justify-content-between align-items-center">
                            <span className="font-bold">Başlangıç:</span>
                            <span>{formatDate(league.startedDate)}</span>
                        </div>
                        <div className="flex justify-content-between align-items-center">
                            <span className="font-bold">Bitiş:</span>
                            <span>{formatDate(league.endDate)}</span>
                        </div>
                    </div>
                </div>

            </div>
        </Card>
    );
}