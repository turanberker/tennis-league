import { useEffect, useState, useRef, useMemo } from "react";
import { Card } from "primereact/card";
import { InputText } from "primereact/inputtext";
import { Button } from "primereact/button";
import { Toast } from "primereact/toast";
import { createFixture, getLeagues, saveLeague } from "../api/leagueService";
import * as yup from "yup";
import { Controller, FormProvider, useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import { classNames } from "primereact/utils";
import { useNavigate } from "react-router-dom";
import { DataTable } from "primereact/datatable";
import { Column } from "primereact/column";
import { Sidebar } from "primereact/sidebar";
import {
  LEAGUE_CATEGORY,
  League_Category_Labels,
  League_Category_Options,
  LEAGUE_FORMAT,
  League_Format_Labels,
  League_Format_Options,
  LEAGUE_PROCESS_TYPE,
  League_Process_Type_Labels,
  League_Process_Type_Options,
  LEAGUE_STATUS,
  League_Status_Labels,
  LeagueListResponse,
  PersistLeagueRequest,
} from "../model/league.model";
import FormItem from "../components/FormItem";
import { Dropdown } from "primereact/dropdown";
import { SplitButton } from "primereact/splitbutton";
import ProtectedRoute from "../router/ProtectedRoute";
import { Role } from "../model/user.model";
import Guard from "../helper/Guard";
import { MenuItem } from "primereact/menuitem";
import { useAuth } from "../context/AuthContext";


// ================= VALIDATION SCHEMA =================
const schema = yup.object({
  name: yup
    .string()
    .required("Lig adı zorunludur.")
    .min(3, "Lig adı en az 3 karakter olmalıdır.")
    .max(75, "Lig adı en fazla 75 karakter olabilir."),
  format: yup.mixed<LEAGUE_FORMAT>().required("Lig formatı zorunludur."),
  category: yup.mixed<LEAGUE_CATEGORY>().required("Lig kategorisi zorunludur."),
  processType: yup
    .mixed<LEAGUE_PROCESS_TYPE>()
    .required("Lig süreç tipi zorunludur."),
});

export default function Leagues() {
  const navigate = useNavigate();
  const { user } = useAuth()
  const [leagues, setLeagues] = useState<LeagueListResponse[]>([]);
  const [selectedLeague, setSelectedLeague] = useState<LeagueListResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [createVisible, setCreateVisible] = useState<boolean>(false);
  const toast = useRef<Toast>(null);

  // ================= REACT HOOK FORM =================
  const methods = useForm<PersistLeagueRequest>({
    resolver: yupResolver(schema),
    defaultValues: {
      name: "",
      format: LEAGUE_FORMAT.Double,
      category: LEAGUE_CATEGORY.MIX,
      processType: LEAGUE_PROCESS_TYPE.FIXTURE,
    },
  });

  const loadLeagues = async () => {
    setLoading(true);
    // Hata olursa 'data' null gelecek ve alt satırlar patlamayacak
    const data = await getLeagues();
    if (data) {
      setLeagues(data);
    }
    setLoading(false); // Hata ol
  };

  useEffect(() => {
    loadLeagues();
  }, []);

  const items: MenuItem[] = useMemo(() => {
    // 1. Lig seçili mi?
    const isSelected = !!selectedLeague;

    // 2. Kullanıcı bu ligin koordinatörü mü? 
    // (Rol kontrolü üstte yapıldığına göre burada sadece ID check yeterli)
    const isCoordinator = user && selectedLeague?.coordinatorUserIds?.includes(user.userID);

    return [
      {
        label: "Fikstür Oluştur",
        icon: "pi pi-plus-circle",
        // ŞART: Seçili değilse VEYA katılım yoksa VEYA listede YOKSA kapat.
        // !isCoordinator -> "Koordinatör değilse" anlamına gelir.
        disabled: !isSelected || selectedLeague?.totalAttentance === 0 || !isCoordinator,
        command: () => handleCreateFixture()
      },
      {
        label: "Takımlar & Oyuncular",
        icon: "pi pi-users",
        disabled: !isSelected || !isCoordinator,
        command: () => handleTeams()
      }
    ];
  }, [selectedLeague, user]);

  const header = () => {
    return (
      <div className="flex justify-content-end">

        <Guard allowedRoles={[Role.ADMIN, Role.COORDINATOR]}>
          <SplitButton label="Yeni Lig" icon="pi pi-plus" size="small" onClick={() => setCreateVisible(true)} model={items} />
        </Guard>
        <Button
          disabled={!selectedLeague}
          rounded
          text
          label="Takımlar & Oyuncular"
          icon="pi pi-chart-bar"
          outlined
          size="small"
          onClick={() => handleTeams()}
        />
        <Button
          disabled={!selectedLeague || selectedLeague.status !== LEAGUE_STATUS.DRAFT}
          rounded
          text
          label="Fikstürü Gör"
          icon="pi pi-calendar"
          outlined
          size="small"
          onClick={() => handleFixtures()}
        />
        <Button
          rounded
          disabled={!selectedLeague || selectedLeague.status !== LEAGUE_STATUS.DRAFT}
          text
          label="Puan Durumu"
          icon="pi pi-chart-bar"
          outlined
          size="small"
          onClick={() => handleStandings()}
        />

      </div>
    );
  };

  const handleStandings = () => {
    navigate(`/leagues/${selectedLeague!.id}/standings`);
  };

  const handleFixtures = () => {
    navigate(`/leagues/${selectedLeague!.id}/fixtures`);
  };

  const handleTeams = () => {
    navigate(`/leagues/${selectedLeague!.id}/teams`);
  };

  // Şemayı ve alan adını verince true/false döner
  const isFieldRequired = (schema: any, fieldName: string) => {
    return schema
      .describe()
      .fields[fieldName]?.tests.some((test: any) => test.name === "required");
  };

  const handleCreateFixture = async () => {
    const data = await createFixture(selectedLeague!.id);
    if (data) {
      toast.current?.show({
        severity: "success",
        summary: "Başarılı",
        detail: "Fikstür başarıyla oluşturuldu",
        life: 3000,
      });
      loadLeagues();
    }
  };

  const onSubmit = async (data: PersistLeagueRequest) => {
    let leagueId = await saveLeague(data);
    if (leagueId) {
      methods.reset();
      toast.current?.show({
        severity: "success",
        summary: "Başarılı",
        detail: "Lig başarıyla oluşturuldu",
        life: 3000,
      });

      setCreateVisible(false);
      loadLeagues(); // listeyi yenile
    }
  };

  return (
    <>
      <Toast ref={toast} />
      <Card
        title="Ligler"
        subTitle="Mevcut ligleri görüntüleyebilir veya yeni lig tanımlayabilirsiniz."
      >
        <DataTable
          value={leagues}
          header={header}
          key="id"
          emptyMessage="Lig bulunamadı"
          loading={loading}
          tableStyle={{ minWidth: "50rem" }}
          dataKey="id"
          selectionMode="single"
          selection={selectedLeague!}
          onSelectionChange={(e) => setSelectedLeague(e.value)}
        >
          <Column
            selectionMode="single"
            headerStyle={{ width: "3rem" }}
          ></Column>
          <Column field="name" header="Lig Adı" />
          <Column field="totalAttentance" header="Katılımcı Sayısı" />
          <Column
            body={(league: LeagueListResponse) => League_Format_Labels[league.format]}
            header="Format"
          />
          <Column
            body={(league: LeagueListResponse) => League_Category_Labels[league.category]}
            header="Category"
          />
          <Column
            body={(league: LeagueListResponse) => League_Process_Type_Labels[league.processType]}
            header="İlerleyiş"
          />
          <Column
            body={(league: LeagueListResponse) => League_Status_Labels[league.status]}
            header="Durumu"
          />
        </DataTable>
      </Card>

      <Sidebar
        header="Yeni Lig Tanımla"
        visible={createVisible}
        position="right"
        onHide={() => setCreateVisible(false)}
        className="w-full md:w-25rem"
      >
        {/* FormProvider ile FormItem'ların context'e erişmesini sağlıyoruz */}
        <FormProvider {...methods}>
          <form onSubmit={methods.handleSubmit(onSubmit)} className="p-fluid">
            <FormItem
              label="Lig Adı"
              name="name"
              required={isFieldRequired(schema, "name")}
            >
              <InputText
                id="name"
                {...methods.register("name")}
                className={classNames({
                  "p-invalid": methods.formState.errors.name,
                })}
                placeholder="Örn: Süper Lig"
              />
            </FormItem>

            <FormItem
              label="Formatı"
              name="format"
              required={isFieldRequired(schema, "format")}
            >
              <Controller
                name="format"
                control={methods.control}
                render={({ field }) => (
                  <Dropdown
                    id="format"
                    {...field}
                    options={League_Format_Options}
                    placeholder="Format Seçin"
                    className={classNames({
                      "p-invalid": methods.formState.errors.format,
                    })}
                  />
                )}
              />
            </FormItem>

            <FormItem
              label="Kategori"
              name="category"
              required={isFieldRequired(schema, "category")}
            >
              <Controller
                name="category"
                control={methods.control}
                render={({ field }) => (
                  <Dropdown
                    id="category"
                    {...field}
                    options={League_Category_Options}
                    placeholder="Kategori Seçin"
                    className={classNames({
                      "p-invalid": methods.formState.errors.category,
                    })}
                  />
                )}
              />
            </FormItem>

            <FormItem
              label="Lig Türü"
              name="processType"
              required={isFieldRequired(schema, "processType")}
            >
              <Controller
                name="processType"
                control={methods.control}
                render={({ field }) => (
                  <Dropdown
                    id="format"
                    {...field}
                    options={League_Process_Type_Options}
                    placeholder="Lig Türü Seçin"
                    className={classNames({
                      "p-invalid": methods.formState.errors.processType,
                    })}
                  />
                )}
              />
            </FormItem>

            {/* FOOTER BUTON ALANI */}
            <div className="mt-4 flex gap-2">
              <Button
                type="button"
                label="İptal"
                icon="pi pi-times"
                outlined
                severity="secondary"
                onClick={() => setCreateVisible(false)}
                className="w-full"
              />
              <Button
                type="submit"
                label="Kaydet"
                icon="pi pi-check"
                loading={methods.formState.isSubmitting}
                className="w-full"
              />
            </div>
          </form>
        </FormProvider>
      </Sidebar>
    </>
  );
}
