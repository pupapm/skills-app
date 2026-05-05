<template>
  <AppShell
    title="Submit Performance"
    subtitle="Provide your yearly role-based data for evaluation. Accuracy is important."
  >
    <div class="card form-card">
      <!-- BASIC -->
      <div class="section">
        <div class="section-head">
          <div>
            <div class="section-kicker">BASIC</div>
            <h3>General Information</h3>
          </div>
        </div>

        <div class="grid-2">
          <label>
            <span>Year</span>
            <input v-model="periodY" placeholder="2026" />
          </label>

          <label>
            <span>Role</span>
            <input :value="session.role || '-'" disabled />
          </label>
        </div>
      </div>

      <!-- UX FORM -->
      <template v-if="session.role === 'ux'">
        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">EXPERIENCE</div>
              <h3>Experience Context</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              ระยะเวลาประสบการณ์ในบทบาท UX/UI เพื่อใช้เป็น context ของความเชี่ยวชาญ
            </div>

            <label>
              <span class="field-title">Experience</span>
              <select v-model="ux.ux_tenure_bucket">
                <option>&lt; 6 เดือน</option>
                <option>6 - 12 เดือน</option>
                <option>1 - 2 ปี</option>
                <option>2 - 4 ปี</option>
                <option>&gt; 4 ปี</option>
              </select>
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">PROJECT OWNERSHIP</div>
              <h3>Scope & Ownership</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน <strong>“โปรเจกต์ที่เข้าไปทำหลัก”</strong> ตลอดปี {{ periodY || "2026" }}

              <div class="hint-subtitle">เช่น:</div>
              <ul>
                <li>คุณเป็นคนคุม direction / decision ของ UX เช่น flow, approach, solution หลัก</li>
                <li>ทีม dev / PM / BA ถามเรื่อง design หลักจะมาหาคุณ</li>
                <li>คุณรับผิดชอบตั้งแต่เริ่ม → handoff</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Main Projects</span>
              <input v-model.number="ux.ux_main_projects" type="number" min="0" />
            </label>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน <strong>“โปรเจกต์ที่เข้าไปช่วย”</strong> ตลอดปี {{ periodY || "2026" }}

              <ul>
                <li>ช่วย review / ช่วยคิดบางส่วนของ flow / ทำบางหน้าจอ</li>
                <li>มีส่วนช่วยให้งานเดิน แต่ไม่ได้คุม direction หลัก</li>
                <li>ไม่ใช่เจ้าของ decision หลักของ project</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Support Projects</span>
              <input v-model.number="ux.ux_support_projects" type="number" min="0" />
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">QUALITY</div>
              <h3>Design Quality & Stability</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวนงานที่ถูก <strong>“ตีกลับ / ต้องแก้ใหม่”</strong> เพราะขาด flow / logic / structure

              <ul>
                <li>flow ไม่ครบ / user journey ขาด</li>
                <li>information ไม่ชัด / edge case ไม่คิด</li>
                <li>งานต้องกลับมาเติมใหม่เพราะ UX scope ไม่ชัด</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Rework Due Scope</span>
              <input v-model.number="ux.rework_due_scope_main" type="number" min="0" />
            </label>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวนโปรเจกต์ที่คุณเสนอ <strong>“มากกว่า 1 solution”</strong> พร้อม reasoning

              <ul>
                <li>มี option A/B พร้อม trade-off เช่น tech / time / risk / user impact</li>
                <li>อธิบายเหตุผลของแต่ละทางเลือกได้</li>
                <li>ไม่ใช่แค่เปลี่ยนสี / layout แบบไม่มีเหตุผล</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Solutions Proposed</span>
              <input v-model.number="ux.solutions_proposed_main" type="number" min="0" />
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">DELIVERY</div>
              <h3>Execution Discipline</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวนโปรเจกต์ที่ dev <strong>“เริ่มทำก่อน UX เข้า”</strong>

              <ul>
                <li>dev เริ่มจาก assumption แล้ว UX มาแก้ทีหลัง</li>
                <li>ทำให้เกิด rework หรือ misalignment</li>
                <li>UX เข้าไปช้าเกินกว่าจะ influence direction ได้เต็มที่</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Late Entry</span>
              <input v-model.number="ux.ux_late_entry_main" type="number" min="0" />
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">IMPACT</div>
              <h3>Post-Release Usability</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน <strong>issue usability ที่เกิดหลังใช้งานจริง</strong> ตลอดปี {{ periodY || "2026" }}

              <div class="hint-subtitle">Includes:</div>
              <ul>
                <li>confusing flow / navigation</li>
                <li>drop-off จาก UX</li>
                <li>user error จาก UI ที่ไม่ชัด</li>
                <li>ปัญหาจาก design / flow / communication</li>
              </ul>

              <div class="hint-subtitle">Do NOT include:</div>
              <ul>
                <li>UI polish เล็ก ๆ</li>
                <li>dev bug ที่ไม่เกี่ยวกับ UX / flow / UI communication</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Usability Issues After Release</span>
              <input v-model.number="ux.usability_issues_after_release" type="number" min="0" />
            </label>
          </div>
        </div>
      </template>

      <!-- QA FORM -->
      <template v-else-if="session.role === 'qa'">
        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">EXPERIENCE</div>
              <h3>Experience Context</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              ระยะเวลาประสบการณ์ในบทบาท QA / Tester เพื่อใช้เป็น context ของความเชี่ยวชาญ
            </div>

            <label>
              <span class="field-title">Experience</span>
              <select v-model="qa.qa_tenure_bucket">
                <option>&lt; 6 เดือน</option>
                <option>6 - 12 เดือน</option>
                <option>1 - 2 ปี</option>
                <option>2 - 4 ปี</option>
                <option>&gt; 4 ปี</option>
              </select>
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">PROJECT OWNERSHIP</div>
              <h3>QA Project Involvement</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน <strong>โปรเจกต์ที่คุณเป็น QA หลัก</strong> ตลอดปี {{ periodY || "2026" }}

              <div class="hint-subtitle">เช่น:</div>
              <ul>
                <li>คุณรับผิดชอบ test strategy / test execution หลัก</li>
                <li>ทีม dev / PM ถามเรื่อง testing status หรือ defect หลักจากคุณ</li>
                <li>คุณดูแลตั้งแต่ test planning → release verification</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Main Projects</span>
              <input v-model.number="qa.qa_main_projects" type="number" min="0" />
            </label>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน <strong>โปรเจกต์ที่คุณเข้าไปช่วย QA</strong> แต่ไม่ได้เป็น owner หลัก

              <ul>
                <li>ช่วย test บาง module</li>
                <li>ช่วย regression / verification</li>
                <li>ช่วย report หรือ review defect บางส่วน</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Support Projects</span>
              <input v-model.number="qa.qa_support_projects" type="number" min="0" />
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">RISK</div>
              <h3>Release Risk & Incidents</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน incident ระดับ critical ที่เกิดขึ้นหลัง release ในโปรเจกต์ที่คุณดูแลหลัก

              <ul>
                <li>production issue ที่กระทบ user / business flow สำคัญ</li>
                <li>issue ที่ควรตรวจพบได้ใน QA process</li>
                <li>ไม่รวม incident ที่เกิดจาก infrastructure หรือ external dependency โดยตรง</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Critical Incident Count</span>
              <input v-model.number="qa.critical_incident_count" type="number" min="0" />
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">COMMUNICATION</div>
              <h3>Defect Communication Quality</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              สัดส่วน bug ที่ dev สามารถแก้ได้โดยไม่ต้องถามข้อมูลเพิ่มจาก QA

              <ul>
                <li>bug report มี step / expected / actual ชัดเจน</li>
                <li>มี evidence เพียงพอ เช่น screenshot, log, payload</li>
                <li>dev เข้าใจปัญหาและเริ่มแก้ได้ทันที</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Fix No Question Ratio</span>
              <select v-model="qa.fix_no_question_ratio_bucket">
                <option>0-20%</option>
                <option>21-40%</option>
                <option>41-60%</option>
                <option>61-80%</option>
                <option>81-100%</option>
              </select>
            </label>
          </div>

          <div class="form-group">
            <div class="field-hint">
              สัดส่วน bug ที่ถูก reopen เพราะข้อมูลจาก QA ไม่พอ / reproduce ไม่ได้ / evidence ไม่ชัด

              <ul>
                <li>dev ต้องถามซ้ำหลายรอบ</li>
                <li>step ไม่ชัด หรือ environment ไม่ครบ</li>
                <li>bug ถูกปิดแล้วกลับมาเปิดใหม่เพราะข้อมูลไม่พอ</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Reopen Due Info Ratio</span>
              <select v-model="qa.reopen_due_info_ratio_bucket">
                <option>0-20%</option>
                <option>21-40%</option>
                <option>41-60%</option>
                <option>61-80%</option>
                <option>81-100%</option>
              </select>
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">TOOLING</div>
              <h3>Evidence & Automation</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน bug ทั้งหมดที่คุณ report ในปี {{ periodY || "2026" }}

              <ul>
                <li>นับเฉพาะ bug ที่ถูกบันทึกในระบบงานจริง เช่น Jira / issue tracker</li>
                <li>ไม่รวม comment เล็ก ๆ ที่ไม่ได้เปิดเป็น bug</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Bugs Reported Total</span>
              <input v-model.number="qa.bugs_reported_total" type="number" min="0" />
            </label>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน bug ที่มี evidence ชัดเจนประกอบ

              <ul>
                <li>screenshot / screen recording</li>
                <li>log / request payload / response</li>
                <li>environment, account, test data หรือ reproduce steps ครบ</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Bug With Evidence Count</span>
              <input v-model.number="qa.bug_with_evidence_count" type="number" min="0" />
            </label>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน automation test หรือ automated check ที่คุณ execute / maintain ในปีนี้

              <ul>
                <li>regression automation</li>
                <li>API test collection</li>
                <li>automated smoke / sanity test</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Automation Executed Count</span>
              <input v-model.number="qa.automation_executed_count" type="number" min="0" />
            </label>
          </div>
        </div>
      </template>

      <!-- BA FORM -->
      <template v-else-if="session.role === 'ba'">
        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">EXPERIENCE</div>
              <h3>Experience Context</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              ระยะเวลาประสบการณ์ในบทบาท Business Analyst เพื่อใช้เป็น context ของความเชี่ยวชาญ
            </div>

            <label>
              <span class="field-title">Experience</span>
              <select v-model="ba.ba_tenure_bucket">
                <option>&lt; 6 เดือน</option>
                <option>6 - 12 เดือน</option>
                <option>1 - 2 ปี</option>
                <option>2 - 4 ปี</option>
                <option>&gt; 4 ปี</option>
              </select>
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">PROJECT OWNERSHIP</div>
              <h3>BA Project Involvement</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน <strong>โปรเจกต์ที่คุณเป็น BA หลัก</strong> ตลอดปี {{ periodY || "2026" }}

              <ul>
                <li>คุณเป็นคนคุม requirement direction หลัก</li>
                <li>stakeholder / dev / QA ถาม requirement หลักจากคุณ</li>
                <li>คุณดูแลตั้งแต่ requirement gathering → signoff / UAT support</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Main Projects</span>
              <input v-model.number="ba.ba_main_projects" type="number" min="0" />
            </label>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน <strong>โปรเจกต์ที่คุณเข้าไปช่วย BA</strong> แต่ไม่ได้เป็น owner หลัก

              <ul>
                <li>ช่วยเก็บ requirement บางส่วน</li>
                <li>ช่วย review document / flow / UAT case</li>
                <li>ช่วยประสานงานบางช่วงของ project</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Support Projects</span>
              <input v-model.number="ba.ba_support_projects" type="number" min="0" />
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">SCOPE</div>
              <h3>Requirement Stability</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวนโปรเจกต์ที่ต้อง rework เพราะ requirement ambiguous / ตีความไม่ตรง / ขาด decision สำคัญ

              <ul>
                <li>dev ต้องกลับมาถาม requirement เดิมซ้ำ</li>
                <li>scope เปลี่ยนเพราะ requirement ไม่ชัดตั้งแต่แรก</li>
                <li>acceptance criteria ไม่พอ ทำให้งานต้องแก้ใหม่</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Req Rework Due Ambiguity</span>
              <input v-model.number="ba.req_rework_due_ambiguity_main" type="number" min="0" />
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">DECISION</div>
              <h3>Requirement Clarity & Signoff</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวนโปรเจกต์ที่ requirement/action item ชัดเจนพอ และไม่ต้องถามซ้ำเรื่องเดิม

              <ul>
                <li>requirement มี action, owner, rule, expected result ชัด</li>
                <li>dev / QA เอาไปทำต่อได้โดยไม่ต้องถามซ้ำหลายรอบ</li>
                <li>ลด repeated clarification ในประเด็นเดิม</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Actionable Req No Repeat</span>
              <input v-model.number="ba.actionable_req_no_repeat_main" type="number" min="0" />
            </label>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวนโปรเจกต์ที่มี decision / signoff ชัดเจนก่อนเริ่ม build

              <ul>
                <li>stakeholder เห็นชอบ requirement ก่อน dev เริ่ม</li>
                <li>มี documented decision หรือ approval ชัด</li>
                <li>ลดการเปลี่ยน requirement ระหว่าง build</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Decision Signoff Before Build</span>
              <input v-model.number="ba.decision_signoff_before_build_main" type="number" min="0" />
            </label>
          </div>
        </div>

        <div class="section">
          <div class="section-head">
            <div>
              <div class="section-kicker">RISK</div>
              <h3>Risk Prevention</h3>
            </div>
          </div>

          <div class="form-group">
            <div class="field-hint">
              จำนวน risk ที่คุณ identify / prevent ก่อนกลายเป็นปัญหาใน project

              <ul>
                <li>เจอ requirement gap ก่อน dev เริ่ม</li>
                <li>จับ conflict ระหว่าง stakeholder ได้ก่อนเกิด rework</li>
                <li>ป้องกัน scope / compliance / process issue ก่อน release</li>
              </ul>
            </div>

            <label>
              <span class="field-title">Risk Prevented Count</span>
              <input v-model.number="ba.risk_prevented_count_main" type="number" min="0" />
            </label>
          </div>
        </div>
      </template>

      <!-- DEV FORM -->
<template v-else-if="session.role === 'dev'">
  <div class="section">
    <div class="section-head">
      <div>
        <div class="section-kicker">EXPERIENCE</div>
        <h3>Experience Context</h3>
      </div>
    </div>

    <div class="form-group">
      <div class="field-hint">
        ระยะเวลาประสบการณ์ในบทบาท Developer เพื่อใช้เป็น context ของความเชี่ยวชาญ
      </div>

      <label>
        <span class="field-title">Experience</span>
        <select v-model="dev.dev_tenure_bucket">
          <option>&lt; 6 เดือน</option>
          <option>6 - 12 เดือน</option>
          <option>1 - 2 ปี</option>
          <option>2 - 4 ปี</option>
          <option>&gt; 4 ปี</option>
        </select>
      </label>
    </div>
  </div>

  <div class="section">
    <div class="section-head">
      <div>
        <div class="section-kicker">PROJECT OWNERSHIP</div>
        <h3>Development Scope & Ownership</h3>
      </div>
    </div>

    <div class="form-group">
      <div class="field-hint">
        จำนวน <strong>โปรเจกต์ที่คุณเป็นคนพัฒนาหลัก</strong> ตลอดปี {{ periodY || "2026" }}

        <ul>
          <li>คุณเป็นคน implement feature หลัก / module หลัก</li>
          <li>เป็นคนแก้ logic สำคัญ หรือ core flow ของระบบ</li>
          <li>ทีม rely on คุณในส่วน dev ของงานนั้น</li>
        </ul>
      </div>

      <label>
        <span class="field-title">Main Projects</span>
        <input v-model.number="dev.dev_main_projects" type="number" min="0" />
      </label>
    </div>

    <div class="form-group">
      <div class="field-hint">
        จำนวน <strong>โปรเจกต์ที่คุณเข้าไปช่วยบางส่วน</strong> ตลอดปี {{ periodY || "2026" }}

        <ul>
          <li>ช่วยแก้ bug / refactor / support feature บางจุด</li>
          <li>ไม่ได้เป็น owner หลักของ module</li>
        </ul>
      </div>

      <label>
        <span class="field-title">Support Projects</span>
        <input v-model.number="dev.dev_support_projects" type="number" min="0" />
      </label>
    </div>
  </div>

  <div class="section">
    <div class="section-head">
      <div>
        <div class="section-kicker">CODE QUALITY</div>
        <h3>Quality & Stability</h3>
      </div>
    </div>

    <div class="form-group">
      <div class="field-hint">
        จำนวน bug ที่เกิดจาก code ของคุณเองหลัง deploy

        <ul>
          <li>logic error</li>
          <li>edge case ที่พลาด</li>
          <li>regression จาก code change</li>
        </ul>
      </div>

      <label>
        <span class="field-title">Bug From Own Code</span>
        <input v-model.number="dev.bug_from_own_code" type="number" min="0" />
      </label>
    </div>

    <div class="form-group">
      <div class="field-hint">
        จำนวนโปรเจกต์ที่ code ของคุณ maintainable และไม่ต้อง rewrite ซ้ำ

        <ul>
          <li>อ่านง่ายและต่อยอดได้</li>
          <li>dev คนอื่นสามารถเข้าใจและแก้ต่อได้</li>
          <li>ไม่ต้อง refactor ใหญ่หลัง merge</li>
        </ul>
      </div>

      <label>
        <span class="field-title">Clean Implementation</span>
        <input v-model.number="dev.clean_implementation" type="number" min="0" />
      </label>
    </div>
  </div>

  <div class="section">
    <div class="section-head">
      <div>
        <div class="section-kicker">DELIVERY</div>
        <h3>Execution Discipline</h3>
      </div>
    </div>

    <div class="form-group">
      <div class="field-hint">
        จำนวนโปรเจกต์ที่คุณส่งงาน dev ตรงเวลา

        <ul>
          <li>ส่งตาม timeline ที่ตกลงกับทีม</li>
          <li>ไม่ delay flow ของ QA / BA / PM / dev คนอื่น</li>
        </ul>
      </div>

      <label>
        <span class="field-title">On-Time Delivery</span>
        <input v-model.number="dev.on_time_delivery" type="number" min="0" />
      </label>
    </div>

    <div class="form-group">
      <div class="field-hint">
        จำนวนครั้งที่งานของคุณกลายเป็น blocker ให้ทีม

        <ul>
          <li>code ไม่เสร็จ ทำให้ QA test ไม่ได้</li>
          <li>dependency ไม่พร้อม ทำให้ทีมอื่นติด</li>
          <li>delay จาก dev task ของคุณกระทบงานคนอื่น</li>
        </ul>
      </div>

      <label>
        <span class="field-title">Blocker Caused</span>
        <input v-model.number="dev.blocker_caused" type="number" min="0" />
      </label>
    </div>
  </div>

  <div class="section">
    <div class="section-head">
      <div>
        <div class="section-kicker">ENGINEERING PRACTICE</div>
        <h3>Technical Maturity</h3>
      </div>
    </div>

    <div class="form-group">
      <div class="field-hint">
        จำนวนครั้งที่คุณช่วย review code ให้ทีม

        <ul>
          <li>ให้ feedback ที่ improve quality</li>
          <li>detect bug ก่อน merge</li>
          <li>ช่วยให้ codebase maintainable ขึ้น</li>
        </ul>
      </div>

      <label>
        <span class="field-title">Code Review Contribution</span>
        <input v-model.number="dev.code_review_contribution" type="number" min="0" />
      </label>
    </div>

    <div class="form-group">
      <div class="field-hint">
        จำนวนงาน optimization / refactor ที่คุณทำ

        <ul>
          <li>improve performance</li>
          <li>refactor code ให้ maintain ง่ายขึ้น</li>
          <li>ลด complexity ของระบบ</li>
        </ul>
      </div>

      <label>
        <span class="field-title">Optimization / Refactor Work</span>
        <input v-model.number="dev.optimization_work" type="number" min="0" />
      </label>
    </div>
  </div>
</template>

      <div class="submit-row">
        <button type="button" @click="submit" :disabled="loading">
          {{ loading ? "Submitting..." : "Submit Data" }}
        </button>
      </div>

      <div v-if="result" class="result-card">
        <div>
          <div class="result-label">Final Score</div>
          <div class="result-value">{{ Number(result.skill_total).toFixed(2) }}</div>
        </div>
        <div>
          <div class="result-label">Credit</div>
          <div class="result-value">{{ Number(result.credit).toFixed(2) }}</div>
        </div>
      </div>

      <p v-if="success" class="success">Submitted successfully</p>
      <p v-if="error" class="error">{{ error }}</p>
    </div>
  </AppShell>
</template>

<script setup>
import { reactive, ref } from "vue";
import AppShell from "../../layouts/AppShell.vue";
import { api, getSession } from "../../api/client";

const session = getSession();
const periodY = ref(String(new Date().getFullYear()));

const ux = reactive({
  ux_tenure_bucket: "1 - 2 ปี",
  ux_main_projects: 0,
  ux_support_projects: 0,
  rework_due_scope_main: 0,
  solutions_proposed_main: 0,
  ux_late_entry_main: 0,
  usability_issues_after_release: 0,
});

const qa = reactive({
  qa_tenure_bucket: "1 - 2 ปี",
  qa_main_projects: 0,
  qa_support_projects: 0,
  bugs_reported_total: 0,
  critical_incident_count: 0,
  fix_no_question_ratio_bucket: "61-80%",
  reopen_due_info_ratio_bucket: "0-20%",
  bug_with_evidence_count: 0,
  automation_executed_count: 0,
});

const ba = reactive({
  ba_tenure_bucket: "1 - 2 ปี",
  ba_main_projects: 0,
  ba_support_projects: 0,
  req_rework_due_ambiguity_main: 0,
  actionable_req_no_repeat_main: 0,
  decision_signoff_before_build_main: 0,
  risk_prevented_count_main: 0,
});

const dev = reactive({
  dev_tenure_bucket: "1 - 2 ปี",
  dev_main_projects: 0,
  dev_support_projects: 0,
  bug_from_own_code: 0,
  clean_implementation: 0,
  on_time_delivery: 0,
  blocker_caused: 0,
  code_review_contribution: 0,
  optimization_work: 0,
});

const loading = ref(false);
const success = ref(false);
const error = ref("");
const result = ref(null);

async function submit() {
  loading.value = true;
  error.value = "";
  success.value = false;
  result.value = null;

  try {
    if (session.role === "ux") {
        result.value = await api.submitUX({ period_y: periodY.value, ...ux });
    } else if (session.role === "qa") {
        result.value = await api.submitQA({ period_y: periodY.value, ...qa });
    } else if (session.role === "ba") {
        result.value = await api.submitBA({ period_y: periodY.value, ...ba });
    } else if (session.role === "dev") {
        result.value = await api.submitDev({ period_y: periodY.value, ...dev });
    } else {
        throw new Error("Unknown role");
    }

    success.value = true;
  } catch (e) {
    error.value = e.message || "Submission failed";
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 18px;
}

.form-card {
  display: grid;
  gap: 28px;
}

.section {
  display: grid;
  gap: 14px;
}

.section-head {
  border-bottom: 1px solid #eee;
  padding-bottom: 8px;
}

.section-head h3 {
  margin: 4px 0 0;
}

.section-kicker {
  font-size: 11px;
  letter-spacing: 0.1em;
  color: #6b7280;
}

.grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.form-group {
  display: grid;
  gap: 8px;
}

.field-title {
  font-weight: 600;
}

.field-hint {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 12px 14px;
  font-size: 13px;
  color: #374151;
  line-height: 1.5;
}

.hint-subtitle {
  margin-top: 8px;
  font-weight: 600;
}

ul {
  margin: 6px 0 0 16px;
  padding: 0;
}

li {
  margin-bottom: 4px;
}

label {
  display: grid;
  gap: 6px;
}

input,
select,
textarea {
  border: 1px solid #d1d5db;
  padding: 10px 12px;
  border-radius: 10px;
  font: inherit;
  background: white;
}

input:disabled {
  background: #f3f4f6;
  color: #6b7280;
}

.submit-row {
  display: flex;
  justify-content: flex-end;
}

button {
  background: #2563eb;
  color: white;
  padding: 10px 16px;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font: inherit;
}

button:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.result-card {
  display: flex;
  gap: 18px;
  flex-wrap: wrap;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 16px;
  padding: 18px;
}

.result-label {
  font-size: 13px;
  color: #1d4ed8;
  margin-bottom: 4px;
}

.result-value {
  font-size: 28px;
  font-weight: 800;
  color: #111827;
}

.success {
  color: #15803d;
}

.error {
  color: #dc2626;
}

@media (max-width: 900px) {
  .grid-2 {
    grid-template-columns: 1fr;
  }
}
</style>