<script lang="ts">
  import { onMount } from "svelte";

  let count = 0;
  let message = "";
  type values = {
    id: number;
    text: string;
  };
  let browsers: string[] = ["brave", "msedge", "firefox"];
  let profiles: string[] = [
    "Default",
    "Profile 1",
    "Profile 2",
    "Profile 3",
    "Profile 4",
    "Profile 5",
    "Profile 6",
    "Profile 7",
    "Profile 8",
    "Profile 9",
  ];
  let settings = {
    browser: browsers[0],
    profile: profiles[0],
  };

  onMount(() => {
    chrome.storage.sync.get(["browser", "profile"], (data) => {
      console.log("sync.get data:", data);
      settings.browser = data.browser;
      settings.profile = data.profile;
    });
    console.log("test");
  });

  const handleSave = () => {
    console.log("settings:", settings);
    chrome.storage.sync
      .set(settings)
      .then(() => {
        message = "更新しました!";
        setTimeout(() => {
          message = "";
        }, 2000);
      })
      .catch((error) => {
        console.log(error);
        message = `error:${error}`;
        setTimeout(() => {
          message = "";
        }, 2000);
      });
  };

  const sendTestOpen = () => {
    chrome.runtime.sendMessage({
      name: "testOpen",
      browser: settings.browser,
      profile: settings.profile,
    });
  };
</script>

<div class="container">
  <main>
    <div class="py-5 text-center">
      <h2>OpenAltBrowserの設定</h2>
      <p class="lead">説明・・・・・・・・</p>
    </div>
    <div class="col-12">
      <h4 class="mb-3">開くブラウザの設定</h4>
      <div class="row g-2">
        <div class="col-md-6">
          <label for="browser" class="form-label">ブラウザー</label>
          <select bind:value={settings.browser} class="form-select" required>
            {#each browsers as browser}
              <option value={browser}>{browser}</option>
            {/each}
          </select>
          <div class="invalid-feedback">対象のブラウザーを選択してください</div>
        </div>
        <div class="col-md-6">
          <label for="profile" class="form-label">プロファイル</label>
          <div class="input-group">
            <select bind:value={settings.profile} class="form-select" required>
              {#each profiles as profile}
                <option value={profile}>{profile}</option>
              {/each}
            </select>
            <button class="btn btn-secondary" on:click={sendTestOpen}
              >ブラウザー起動テスト</button
            >
          </div>
          <div class="invalid-feedback">
            ブラウザーのプロファイルを選択してください
          </div>
        </div>
      </div>
      <hr class="my-4" />
      <button class="btn btn-primary" on:click={handleSave}>保存</button>
      {#if message != ""}
        <div class="alert alert-primary" role="alert">{message}</div>
      {/if}
    </div>
  </main>
</div>
